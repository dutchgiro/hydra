package load_balancer

import (
	"log"
	"runtime"
	"time"

	zmq "github.com/innotech/hydra/vendors/github.com/alecthomas/gozmq"
)

type Requester interface {
	Close()
	Send([]byte, [][]byte) [][]byte
}

type client struct {
	socket         *zmq.Socket
	context        *zmq.Context
	server         string
	timeout        time.Duration
	requestTimeout time.Duration
}

func NewClient(server string, requestTimeout int) *client {
	context, _ := zmq.NewContext()
	self := &client{
		server:         server,
		context:        context,
		timeout:        2500 * time.Millisecond,
		requestTimeout: time.Duration(requestTimeout) * time.Millisecond,
	}
	self.connect()
	runtime.SetFinalizer(self, closeClient)
	return self
}

func (self *client) connect() {
	if self.socket != nil {
		self.socket.Close()
	}

	var errSocket error
	self.socket, errSocket = self.context.NewSocket(zmq.REQ)
	if errSocket != nil {
		log.Fatal(errSocket.Error())
	}
	self.socket.SetLinger(0)
	self.socket.Connect(self.server)
	if err := self.socket.SetRcvTimeout(self.requestTimeout); err != nil {
		log.Println(err)
	}
}

// Close socket and context connections
func closeClient(self *client) {
	if self.socket != nil {
		self.socket.Close()
	}
	self.context.Close()
}

// Close socket and context connections
func (self *client) Close() {
	if self.socket != nil {
		self.socket.Close()
	}
	self.context.Close()
}

// Send dispatchs requests to load balancer server and returns the response message
func (self *client) Send(service []byte, request [][]byte) (reply [][]byte) {
	frame := append([][]byte{service}, request...)

	self.socket.SendMultipart(frame, zmq.NOBLOCK)
	msg, err := self.socket.RecvMultipart(0)

	if err != nil || len(msg) < 1 {
		reply = [][]byte{}
	} else {
		reply = msg
	}

	return
}
