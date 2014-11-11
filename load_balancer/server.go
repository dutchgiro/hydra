package load_balancer

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	zmq "github.com/innotech/hydra/vendors/github.com/alecthomas/gozmq"

	"github.com/innotech/hydra/log"
	. "github.com/innotech/hydra/model/entity"
)

const (
	INTERNAL_SERVICE_PREFIX = "isb."
	// Merge all heartbeat
	HEARTBEAT_INTERVAL = 2500 * time.Millisecond
	HEARTBEAT_EXPIRY   = HEARTBEAT_INTERVAL * HEARTBEAT_LIVENESS
)

type Broker interface {
	Close()
	Run()
}

type lbWorker struct {
	identity      string    // UUID Identity of worker
	address       []byte    // Address to route to
	expiry        time.Time // Expires at this point, unless heartbeat
	priorityLevel int       // Sets the priority level of a worker (level 0 is reserved for workers deployed on the same host)
	// priority int        // Sets the priority of a worker within their priority level
	service *lbService // Owning service, if known
}

type lbChain struct {
	app      string
	client   string
	msg      []byte
	shackles *ZList
}

type lbShackle struct {
	serviceName string
	serviceArgs map[string]interface{}
}

type lbService struct {
	broker   Broker
	name     string
	requests [][][]byte // List of client requests
	waiting  []*ZList   // Lists of waiting workers sorted by priority level
}

type loadBalancer struct {
	context       *zmq.Context          // Context
	heartbeatAt   time.Time             // When to send HEARTBEAT
	services      map[string]*lbService // Known services
	frontend      *zmq.Socket           // Socket for clients
	inprocBackend *zmq.Socket           // Socket for local workers
	tcpBackend    *zmq.Socket           // Socket for external workers
	waiting       *ZList                // Idle workers
	workers       map[string]*lbWorker  // Known workers
	chains        map[string]lbChain    // Peding Requests
}

func NewLoadBalancer(frontendEndpoint, tcpBackendEndpoint, inprocBackendEndpoint string) *loadBalancer {
	context, _ := zmq.NewContext()
	// Define inproc socket to talk with HTTP CLient API
	frontend, _ := context.NewSocket(zmq.ROUTER)
	frontend.SetLinger(0)
	frontend.Bind(frontendEndpoint)
	// Define inproc socket to talk with local Workers
	inprocBackend, _ := context.NewSocket(zmq.ROUTER)
	inprocBackend.SetLinger(0)
	inprocBackend.Bind(inprocBackendEndpoint)
	// Define tcp socket to talk with external Workers
	tcpBackend, _ := context.NewSocket(zmq.ROUTER)
	tcpBackend.SetLinger(0)
	tcpBackend.Bind(tcpBackendEndpoint)
	return &loadBalancer{
		context:       context,
		heartbeatAt:   time.Now().Add(HEARTBEAT_INTERVAL),
		services:      make(map[string]*lbService),
		frontend:      frontend,
		inprocBackend: inprocBackend,
		tcpBackend:    tcpBackend,
		waiting:       NewList(),
		workers:       make(map[string]*lbWorker),
		chains:        make(map[string]lbChain),
	}
}

// Deletes worker from all data structures, and deletes worker.
func (self *loadBalancer) deleteWorker(worker *lbWorker, disconnect bool) {
	if worker == nil {
		log.Warn("Nil worker")
	}

	if disconnect {
		log.Info("SIGNAL_DISCONNECT")
		self.sendToWorker(worker, SIGNAL_DISCONNECT, nil, nil)
	}

	if worker.service != nil {
		log.Info("Delete Worker Service")
		worker.service.waiting[worker.priorityLevel].Delete(worker)
	}
	self.waiting.Delete(worker)
	delete(self.workers, worker.identity)
}

// decomposeMapOfInstancesMsg extracts the uris from nested array of instances
func (self *loadBalancer) decomposeMapOfInstancesMsg(msg []byte) []byte {
	log.Debug("--- decomposeMapOfInstancesMsg")
	var levels []interface{}
	if err := json.Unmarshal(msg, &levels); err != nil {
		log.Warn("Bad JSON format")
	}
	log.Debug("--- levels", levels)
	var computedInstances []interface{}
	var processMapLevels func([]interface{})
	processMapLevels = func(levels []interface{}) {
		for _, level := range levels {
			if level != nil {
				kind := reflect.TypeOf(level).Kind()
				if kind == reflect.Slice || kind == reflect.Array {
					processMapLevels(level.([]interface{}))
				} else {
					computedInstances = append(computedInstances, levels...)
					return
				}
			}
		}
	}
	processMapLevels(levels)
	// TODO: extract uris
	sortedInstanceUris := make([]string, 0)
	for _, instance := range computedInstances {
		if instance != nil {
			if uri, ok := instance.(map[string]interface{})["Info"].(map[string]interface{})["uri"]; ok {
				sortedInstanceUris = append(sortedInstanceUris, uri.(string))
			}
		} else {
			log.Warn("Remove instance without uri attribute from last balancer response")
		}
	}

	log.Debug("--- sortedInstanceUris", sortedInstanceUris)
	uris, err := json.Marshal(sortedInstanceUris)
	if err != nil {
		log.Warn("Bad JSON format")
	}
	return uris
}

// Dispatch chains advancing a shackle
func (self *loadBalancer) advanceShackle(chain lbChain) {
	log.Debug("--- advanceShackle")
	elem := chain.shackles.Pop()
	if elem == nil {
		log.Debug("--- Decompose")
		// Decompose
		instanceUrisMsg := self.decomposeMapOfInstancesMsg(chain.msg)
		log.Debug("--- instanceUrisMsg:", instanceUrisMsg)
		msg := [][]byte{[]byte(chain.client), nil, instanceUrisMsg}
		self.frontend.SendMultipart(msg, 0)
		return
	}
	log.Debug("--- NO Decompose")
	shackle, _ := elem.Value.(lbShackle)
	args, _ := json.Marshal(shackle.serviceArgs)
	msg := [][]byte{[]byte(chain.client), nil, chain.msg, args}
	log.Debug("--- SENDING TO:", shackle.serviceName)
	self.dispatch(self.requireService(shackle.serviceName), msg)
}

// getFirstLevelOfWaitingWorkers returns the first level of workers of a service that is not empty
func (self *loadBalancer) getFirstLevelOfWaitingWorkers(service *lbService) *ZList {
	for i := 0; i < len(service.waiting); i++ {
		if service.waiting[i].Len() > 0 {
			return service.waiting[i]
		}
	}

	return nil
}

// Dispatch requests to waiting workers as possible
func (self *loadBalancer) dispatch(service *lbService, msg [][]byte) {
	if service == nil {
		log.Warn("Nil service")
		return
	}
	// Queue message if any
	if len(msg) != 0 {
		service.requests = append(service.requests, msg)
	}
	self.purgeWorkers()

	workers := self.getFirstLevelOfWaitingWorkers(service)
	if workers == nil {
		return
	}

	for workers.Len() > 0 && len(service.requests) > 0 {
		msg, service.requests = service.requests[0], service.requests[1:]
		elem := workers.Pop()
		self.waiting.Delete(elem.Value)
		worker, _ := elem.Value.(*lbWorker)

		self.sendToWorker(worker, SIGNAL_REQUEST, nil, msg)
	}
}

// Register chain from new client request
func (self *loadBalancer) registerChain(client []byte, msg [][]byte) {
	// log.Debug("--- registerChain")
	var services []Balancer
	if err := json.Unmarshal(msg[1], &services); err != nil {
		panic(err)
	}
	chain := lbChain{
		app:      string(msg[0]),
		client:   string(client),
		msg:      msg[2],
		shackles: NewList(),
	}
	for _, service := range services {
		args := service.Args
		args["appId"] = chain.app
		chain.shackles.PushBack(lbShackle{
			serviceName: service.Id,
			serviceArgs: args,
		})
	}
	// log.Debug("--- CHAIN:", chain)
	self.chains[string(client)] = chain
}

// Process a request coming from a client.
func (self *loadBalancer) processClient(client []byte, msg [][]byte) {
	// Application + Services + Instances
	if len(msg) < 3 {
		log.Warn("Invalid message from client sender")
		return
	}
	// Register chain
	self.registerChain(client, msg)
	// Start chain of requests
	self.advanceShackle(self.chains[string(client)])
}

// Process message sent to us by a worker.
func (self *loadBalancer) processWorker(sender []byte, msg [][]byte) {
	//  At least, command
	if len(msg) < 1 {
		log.Warn("Invalid message from Worker, this doesn contain command")
		return
	}

	command, msg := msg[0], msg[1:]
	identity := hex.EncodeToString(sender)
	worker, workerReady := self.workers[identity]
	if !workerReady {
		worker = &lbWorker{
			identity: identity,
			address:  sender,
			expiry:   time.Now().Add(HEARTBEAT_EXPIRY),
		}
		self.workers[identity] = worker
		log.Debugf("Worker: %#v", worker)
	}

	switch string(command) {
	case SIGNAL_READY:
		log.Debug("SIGNAL_READY")
		//  At least, a service name
		if len(msg) < 1 {
			log.Warn("Invalid message from worker, service name is missing")
			self.deleteWorker(worker, true)
			return
		}
		service := msg[0]
		priorityLevel := msg[1]
		//  Not first command in session or Reserved service name
		if workerReady || string(service[:4]) == INTERNAL_SERVICE_PREFIX {
			self.deleteWorker(worker, true)
		} else {
			//  Attach worker to service and mark as idle
			worker.service = self.requireService(string(service))
			priorityLevelInt, err := strconv.Atoi(string(priorityLevel))
			if err != nil {
				log.Warn("The priority level of a worker must be an integer")
				self.deleteWorker(worker, true)
			}
			worker.priorityLevel = priorityLevelInt
			for len(worker.service.waiting)-1 < priorityLevelInt {
				worker.service.waiting = append(worker.service.waiting, NewList())
			}
			self.workerWaiting(worker)
		}
	case SIGNAL_REPLY:
		log.Debug("SIGNAL_REPLY")
		if workerReady {
			log.Debug("--- workerReady")
			//  Remove & save client return envelope and insert the
			//  protocol header and service name, then rewrap envelope.
			client := msg[0]
			chain := self.chains[string(client)]
			chain.msg = msg[2]
			self.advanceShackle(chain)
			self.workerWaiting(worker)
		} else {
			log.Debug("--- NO workerReady")
			self.deleteWorker(worker, true)
		}
	case SIGNAL_HEARTBEAT:
		log.Debug("SIGNAL_HEARTBEAT")
		if workerReady {
			worker.expiry = time.Now().Add(HEARTBEAT_EXPIRY)
		} else {
			self.deleteWorker(worker, true)
		}
	case SIGNAL_DISCONNECT:
		log.Debug("SIGNAL_DISCONNECT")
		self.deleteWorker(worker, false)
	default:
		log.Warn("Invalid message in Load Balancer")
		Dump(msg)
	}
	return
}

//  Look for & kill expired workers.
//  Workers are oldest to most recent, so we stop at the first alive worker.
func (self *loadBalancer) purgeWorkers() {
	now := time.Now()
	for elem := self.waiting.Front(); elem != nil; elem = elem.Next() {
		worker, _ := elem.Value.(*lbWorker)
		if worker.expiry.After(now) {
			continue
		}
		self.deleteWorker(worker, false)
	}
}

// Locates the service (creates if necessary).
func (self *loadBalancer) requireService(name string) *lbService {
	if len(name) == 0 {
		log.Warn("Invalid service name have been required")
	}
	service, ok := self.services[name]
	if !ok {
		service = &lbService{
			name: name,
			// waiting: NewList(),
			waiting: make([]*ZList, 0),
		}
		self.services[name] = service
	}
	return service
}

// Send message to worker.
// If message is provided, sends that message.
func (self *loadBalancer) sendToWorker(worker *lbWorker, command string, option []byte, msg [][]byte) {
	//  Stack routing and protocol envelopes to start of message and routing envelope
	if len(option) > 0 {
		msg = append([][]byte{option}, msg...)
	}
	msg = append([][]byte{worker.address, nil, []byte(command)}, msg...)

	if worker.priorityLevel == 0 {
		self.inprocBackend.SendMultipart(msg, 0)
	} else {
		self.tcpBackend.SendMultipart(msg, 0)
	}
}

// This worker is now waiting for work.
func (self *loadBalancer) workerWaiting(worker *lbWorker) {
	//  Queue to broker and service waiting lists
	self.waiting.PushBack(worker)
	worker.service.waiting[worker.priorityLevel].PushBack(worker)
	worker.expiry = time.Now().Add(HEARTBEAT_EXPIRY)
	self.dispatch(worker.service, nil)
}

func (self *loadBalancer) Close() {
	if self.frontend != nil {
		self.frontend.Close()
	}
	if self.tcpBackend != nil {
		self.tcpBackend.Close()
	}
	if self.inprocBackend != nil {
		self.inprocBackend.Close()
	}
	self.context.Close()
}

// Main broker working loop
func (self *loadBalancer) Run() {
	for {
		items := zmq.PollItems{
			zmq.PollItem{Socket: self.frontend, Events: zmq.POLLIN},
			zmq.PollItem{Socket: self.tcpBackend, Events: zmq.POLLIN},
			zmq.PollItem{Socket: self.inprocBackend, Events: zmq.POLLIN},
		}

		_, err := zmq.Poll(items, HEARTBEAT_INTERVAL)
		if err != nil {
			log.Fatal("Non items for polling")
		}

		if items[0].REvents&zmq.POLLIN != 0 {
			msg, _ := self.frontend.RecvMultipart(0)
			// TODO: check msg parts
			requestId := msg[0]
			msg = msg[2:]
			self.processClient(requestId, msg)
		}
		if items[1].REvents&zmq.POLLIN != 0 {
			msg, _ := self.tcpBackend.RecvMultipart(0)
			sender := msg[0]
			msg = msg[2:]
			self.processWorker(sender, msg)
		}
		if items[2].REvents&zmq.POLLIN != 0 {
			msg, _ := self.inprocBackend.RecvMultipart(0)
			sender := msg[0]
			msg = msg[2:]
			self.processWorker(sender, msg)
		}

		if self.heartbeatAt.Before(time.Now()) {
			self.purgeWorkers()
			for elem := self.waiting.Front(); elem != nil; elem = elem.Next() {
				worker, _ := elem.Value.(*lbWorker)
				self.sendToWorker(worker, SIGNAL_HEARTBEAT, nil, nil)
			}
			self.heartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
		}
	}
}
