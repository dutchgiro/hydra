package load_balancer

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
	"runtime"
	"strconv"
	"time"

	zmq "github.com/innotech/hydra/vendors/github.com/pebbe/zmq4"

	"github.com/innotech/hydra/log"
	. "github.com/innotech/hydra/model/entity"
)

// TODO: Maybe move to init.go
const (
	// TODO: Maybe remove
	INTERNAL_SERVICE_PREFIX = "isb."

	HEARTBEAT_LIVENESS = 3                       //  3-5 is reasonable
	HEARTBEAT_INTERVAL = 2500 * time.Millisecond //  msecs
	HEARTBEAT_EXPIRY   = HEARTBEAT_INTERVAL * HEARTBEAT_LIVENESS
)

type Broker struct {
	chains                map[string]Chain // Peding Requests
	context               *zmq.Context     //  Context for broker
	frontendEndpoint      string
	frontendSocket        *zmq.Socket //  Socket for clients
	inprocBackendEndpoint string
	inprocBackendSocket   *zmq.Socket //  Socket for local workers
	tcpBackendEndpoint    string
	tcpBackendSocket      *zmq.Socket         //  Socket for external workers
	verbose               bool                //  Print activity to stdout
	endpoint              string              //  Broker binds to this endpoint
	services              map[string]*Service //  Hash of known services
	workers               map[string]*Worker  //  Hash of known workers
	waiting               *ZList              //  List of waiting workers
	heartbeatAt           time.Time           //  When to send HEARTBEAT
}

type Service struct {
	broker   *Broker    //  Broker instance
	name     string     //  Service name
	requests [][][]byte //  List of client requests
	waiting  []*ZList   // Lists of waiting workers sorted by priority level
}

//  The worker class defines a single worker, idle or active:

// type Worker struct {
// 	broker    *Broker   //  Broker instance
// 	id_string string    //  Identity of worker as string
// 	identity  string    //  Identity frame for routing
// 	service   *Service  //  Owning service, if known
// 	expiry    time.Time //  Expires at unless heartbeat
// }

type Chain struct {
	app          string
	client       string
	clientParams []byte
	msg          []byte
	shackles     *ZList
}

type Shackle struct {
	serviceName string
	serviceArgs map[string]interface{}
}

type Worker struct {
	identity      string    // UUID Identity of worker
	address       []byte    // Address to route to
	expiry        time.Time // Expires at this point, unless heartbeat
	priorityLevel int       // Sets the priority level of a worker (level 0 is reserved for workers deployed on the same host)
	service       *Service  // Owning service, if known
}

func NewBroker(frontendEndpoint, tcpBackendEndpoint, inprocBackendEndpoint string) (broker *Broker) {
	broker = &Broker{
		frontendEndpoint:      frontendEndpoint,
		inprocBackendEndpoint: inprocBackendEndpoint,
		tcpBackendEndpoint:    tcpBackendEndpoint,
		chains:                make(map[string]Chain),
		services:              make(map[string]*Service),
		// services:      make(map[string]*lbService),
		workers: make(map[string]*Worker),
		// workers:       make(map[string]*lbWorker),
		// waiting: make([]*Worker, 0),
		waiting:     NewList(),
		heartbeatAt: time.Now().Add(HEARTBEAT_INTERVAL),
	}

	var err error
	broker.context, err = zmq.NewContext()
	if err != nil {
		log.Fatal("LoadBalancer broker ConnectToBroker() creating context failed")
	}

	// Define inproc socket to talk with HTTP CLient API
	broker.frontendSocket, err = broker.context.NewSocket(zmq.ROUTER)
	if err != nil {
		log.Fatal("LoadBalancer broker ConnectToBroker() creating frontend socket failed")
	}
	broker.frontendSocket.SetRcvhwm(500000)
	// TODO: Maybe set linger to 0

	// Define inproc socket to talk with local Workers
	broker.inprocBackendSocket, err = broker.context.NewSocket(zmq.ROUTER)
	if err != nil {
		log.Fatal("LoadBalancer client ConnectToBroker() creating inproc backend socket failed")
	}
	broker.inprocBackendSocket.SetRcvhwm(500000)
	// TODO: Maybe set linger to 0

	// Define tcp socket to talk with external Workers
	broker.tcpBackendSocket, err = broker.context.NewSocket(zmq.ROUTER)
	if err != nil {
		log.Fatal("LoadBalancer client ConnectToBroker() creating tcp backend socket failed")
	}
	broker.tcpBackendSocket.SetRcvhwm(500000)
	// TODO: Maybe set linger to 0

	runtime.SetFinalizer(broker, (*Broker).Close)
	return
}

func (b *Broker) Close() {
	if b.frontendSocket != nil {
		b.frontendSocket.Close()
		b.frontendSocket = nil
	}
	if b.inprocBackendSocket != nil {
		b.inprocBackendSocket.Close()
		b.inprocBackendSocket = nil
	}
	if b.tcpBackendSocket != nil {
		b.tcpBackendSocket.Close()
		b.tcpBackendSocket = nil
	}
	if b.context != nil {
		b.context.Term()
		b.context = nil
	}
	return
}

// Register chain from new client request
func (b *Broker) registerChain(client []byte, msg [][]byte) {
	var services []Balancer
	_ = json.Unmarshal(msg[1], &services)

	chain := Chain{
		app:          string(msg[0]),
		client:       string(client),
		clientParams: msg[3],
		msg:          msg[2],
		shackles:     NewList(),
	}
	for _, service := range services {
		args := service.Args
		args["appId"] = chain.app
		chain.shackles.PushBack(Shackle{
			serviceName: service.Id,
			serviceArgs: args,
		})
	}
	b.chains[string(client)] = chain
}

// decomposeMapOfInstancesMsg extracts the uris from nested array of instances
func (b *Broker) decomposeMapOfInstancesMsg(msg []byte) []byte {
	var levels []interface{}
	if err := json.Unmarshal(msg, &levels); err != nil {
		log.Warn("LoadBalancer broker can not decompose instances message")
		uris, _ := json.Marshal(make([]string, 0))
		return uris
	}
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

	sortedInstanceUris := make([]string, 0)
	for _, instance := range computedInstances {
		if instance != nil {
			if uri, ok := instance.(map[string]interface{})["Info"].(map[string]interface{})["uri"]; ok {
				sortedInstanceUris = append(sortedInstanceUris, uri.(string))
			}
		} else {
			log.Warn("LoadBalancer broker remove instance without uri attribute from last balancer response")
		}
	}

	uris, _ := json.Marshal(sortedInstanceUris)
	return uris
}

// Locates the service (creates if necessary).
func (b *Broker) requireService(name string) *Service {
	if len(name) == 0 {
		log.Warn("LoadBalancer broker receives an invalid service name have been required")
	}
	service, ok := b.services[name]
	if !ok {
		service = &Service{
			name:    name,
			waiting: make([]*ZList, 0),
		}
		b.services[name] = service
	}
	return service
}

// getFirstLevelOfWaitingWorkers returns the first level of workers of a service that is not empty
func (b *Broker) getFirstLevelOfWaitingWorkers(service *Service) *ZList {
	for i := 0; i < len(service.waiting); i++ {
		if service.waiting[i].Len() > 0 {
			return service.waiting[i]
		}
	}

	return nil
}

// Dispatch requests to waiting workers as possible
func (b *Broker) dispatch(service *Service, msg [][]byte) {
	if service == nil {
		log.Warn("LoadBalancer broker tries to dispatch to a nil service")
	}
	// Queue message if any
	if len(msg) != 0 {
		service.requests = append(service.requests, msg)
	}
	b.purgeWorkers()

	workers := b.getFirstLevelOfWaitingWorkers(service)
	if workers == nil {
		return
	}

	for workers.Len() > 0 && len(service.requests) > 0 {
		msg, service.requests = service.requests[0], service.requests[1:]
		elem := workers.Pop()
		b.waiting.Delete(elem.Value)
		worker, _ := elem.Value.(*Worker)

		b.sendToWorker(worker, SIGNAL_REQUEST, nil, msg)
	}
}

// Dispatch chains advancing a shackle
func (b *Broker) advanceShackle(chain Chain) {
	elem := chain.shackles.Pop()
	if elem == nil {
		log.Debugf("LoadBalancer broker end of chain with last message %q\n", chain.msg)
		instanceUrisMsg := b.decomposeMapOfInstancesMsg(chain.msg)
		msg := [][]byte{[]byte(chain.client), nil, instanceUrisMsg}
		b.frontendSocket.SendMessage(msg)
		return
	}
	shackle, _ := elem.Value.(Shackle)
	args, _ := json.Marshal(shackle.serviceArgs)
	msg := [][]byte{[]byte(chain.client), nil, chain.msg, chain.clientParams, args}
	b.dispatch(b.requireService(shackle.serviceName), msg)
}

// Process a request coming from a client.
func (b *Broker) processClientMsg(client []byte, msg [][]byte) {
	// Application + Services + Instances + Query params of client request
	if len(msg) < 4 {
		log.Warn("LoadBalancer broker invalid message received from client sender")
	}

	b.registerChain(client, msg)
	// Start chain of requests
	b.advanceShackle(b.chains[string(client)])
}

// This worker is now waiting for work.
func (b *Broker) workerWaiting(worker *Worker) {
	// Queue to broker and service waiting lists
	b.waiting.PushBack(worker)
	worker.service.waiting[worker.priorityLevel].PushBack(worker)
	worker.expiry = time.Now().Add(HEARTBEAT_EXPIRY)
	b.dispatch(worker.service, nil)
}

//  processWorkerMsg method processes one READY, REPLY, HEARTBEAT or
//  DISCONNECT message sent to the broker by a worker:
// func (broker *Broker) WorkerMsg(sender string, msg []string) {
func (b *Broker) processWorkerMsg(sender []byte, msg [][]byte) {
	//  At least, command
	if len(msg) == 0 {
		log.Warn("LoadBalancer broker process invalid message from Worker, this doesn't contain command")
		return
	}

	// TODO: make a function
	command, msg := msg[0], msg[1:]
	identity := hex.EncodeToString(sender)
	worker, workerReady := b.workers[identity]
	if !workerReady {
		worker = &Worker{
			identity: identity,
			address:  sender,
			expiry:   time.Now().Add(HEARTBEAT_EXPIRY),
		}
		b.workers[identity] = worker
	}

	switch string(command) {
	case SIGNAL_READY:
		//  At least, a service name
		if len(msg) < 1 {
			log.Warn("Invalid message from worker, service name is missing")
			b.deleteWorker(worker, true)
			return
		}
		service := msg[0]
		priorityLevel := msg[1]
		//  Not first command in session or Reserved service name
		if workerReady || string(service[:4]) == INTERNAL_SERVICE_PREFIX {
			b.deleteWorker(worker, true)
		} else {
			//  Attach worker to service and mark as idle
			worker.service = b.requireService(string(service))
			priorityLevelInt, err := strconv.Atoi(string(priorityLevel))
			if err != nil {
				log.Warn("The priority level of a worker must be an integer")
				b.deleteWorker(worker, true)
				return
			}
			worker.priorityLevel = priorityLevelInt
			for len(worker.service.waiting)-1 < priorityLevelInt {
				worker.service.waiting = append(worker.service.waiting, NewList())
			}
			b.workerWaiting(worker)
		}
	case SIGNAL_REPLY:
		log.Debugf("LoadBalancer broker receive SIGNAL_REPLY with message %q\n", msg)
		if workerReady {
			//  Remove & save client return envelope and insert the
			//  protocol header and service name, then rewrap envelope.
			client := msg[0]
			chain := b.chains[string(client)]
			chain.msg = msg[2]
			b.advanceShackle(chain)
			b.workerWaiting(worker)
		} else {
			b.deleteWorker(worker, true)
		}
	case SIGNAL_HEARTBEAT:
		if workerReady {
			worker.expiry = time.Now().Add(HEARTBEAT_EXPIRY)
		} else {
			b.deleteWorker(worker, true)
		}
	case SIGNAL_DISCONNECT:
		b.deleteWorker(worker, false)
	default:
		log.Warn("Load Balancer broker receives invalid message from worker")
		Dump(msg)
	}
	return
}

// Deletes worker from all data structures, and deletes worker.
func (b *Broker) deleteWorker(worker *Worker, disconnect bool) {
	if worker == nil {
		log.Warn("LoadBalancer broker tries to delete nil Worker")
		return
	}

	if disconnect {
		b.sendToWorker(worker, SIGNAL_DISCONNECT, nil, nil)
	}

	if worker.service != nil {
		log.Debug("LoadBalancer broker deletes Worker service")
		worker.service.waiting[worker.priorityLevel].Delete(worker)
	}
	b.waiting.Delete(worker)
	delete(b.workers, worker.identity)
}

// Look for & kill expired workers.
// Workers are oldest to most recent, so we stop at the first alive worker.
func (b *Broker) purgeWorkers() {
	now := time.Now()
	for elem := b.waiting.Front(); elem != nil; elem = elem.Next() {
		worker, _ := elem.Value.(*Worker)
		if worker.expiry.After(now) {
			continue
		}
		log.Debug("LoadBalancer broker deleting expired worker:", worker.identity)
		b.deleteWorker(worker, false)
	}
}

// Send message to worker.
// If message is provided, sends that message.
func (b *Broker) sendToWorker(worker *Worker, command string, option []byte, msg [][]byte) {
	//  Stack routing and protocol envelopes to start of message and routing envelope
	if len(option) > 0 {
		msg = append([][]byte{option}, msg...)
	}
	msg = append([][]byte{worker.address, nil, []byte(command)}, msg...)

	log.Debugf("LoadBalancer broker sending %q to worker %s\n", msg, worker.identity)
	if worker.priorityLevel == 0 {
		b.inprocBackendSocket.SendMessage(msg)
	} else {
		b.tcpBackendSocket.SendMessage(msg)
	}
}

// Run executes the working loop
func (b *Broker) Run() {
	var err error
	// Bind frontend socket
	err = b.frontendSocket.Bind(b.frontendEndpoint)
	if err != nil {
		log.Fatal("LoadBalancer broker failed to bind frontend socket at", b.frontendEndpoint)
	}
	log.Info("LoadBalancer broker frontend socket is active at", b.frontendEndpoint)

	// Bind inproc backend socket
	err = b.inprocBackendSocket.Bind(b.inprocBackendEndpoint)
	if err != nil {
		log.Fatal("LoadBalancer broker failed to bind frontend socket at", b.inprocBackendEndpoint)
	}
	log.Info("LoadBalancer broker inproc backend socket is active at", b.inprocBackendEndpoint)

	// Bind tcp backend socket
	err = b.tcpBackendSocket.Bind(b.tcpBackendEndpoint)
	if err != nil {
		log.Fatal("LoadBalancer broker failed to bind frontend socket at", b.tcpBackendEndpoint)
	}
	log.Info("LoadBalancer broker tcp backend socket is active at", b.tcpBackendEndpoint)

	poller := zmq.NewPoller()
	poller.Add(b.frontendSocket, zmq.POLLIN)
	poller.Add(b.inprocBackendSocket, zmq.POLLIN)
	poller.Add(b.tcpBackendSocket, zmq.POLLIN)

	for {
		polled, err := poller.Poll(HEARTBEAT_INTERVAL)
		if err != nil {
			log.Fatal("LoadBalancer broker non items for polling")
		}

		for _, item := range polled {
			switch socket := item.Socket; socket {
			case b.frontendSocket:
				msg, err := b.frontendSocket.RecvMessageBytes(0)
				if err != nil {
					continue
				}
				log.Debugf("LoadBalancer broker frontend socket received message: %q\n", msg)

				// TODO: check msg parts and extract to function
				requestId := msg[0]
				msg = msg[2:]
				b.processClientMsg(requestId, msg)
			case b.inprocBackendSocket:
				msg, err := b.inprocBackendSocket.RecvMessageBytes(0)
				if err != nil {
					continue
				}
				log.Debugf("LoadBalancer broker inproc backend socket received message: %q\n", msg)

				// TODO: check msg parts and extract to function
				sender := msg[0]
				msg = msg[2:]
				b.processWorkerMsg(sender, msg)
			case b.tcpBackendSocket:
				msg, err := b.tcpBackendSocket.RecvMessageBytes(0)
				if err != nil {
					continue
				}
				log.Debugf("LoadBalancer broker tcp backend socket received message: %q\n", msg)

				// TODO: check msg parts and extract to function
				sender := msg[0]
				msg = msg[2:]
				b.processWorkerMsg(sender, msg)
			}
		}

		if b.heartbeatAt.Before(time.Now()) {
			b.purgeWorkers()
			for elem := b.waiting.Front(); elem != nil; elem = elem.Next() {
				worker, _ := elem.Value.(*Worker)
				b.sendToWorker(worker, SIGNAL_HEARTBEAT, nil, nil)
			}
			b.heartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
		}
	}
	log.Info("LoadBalancer broker interrupt received, shutting down...")
}
