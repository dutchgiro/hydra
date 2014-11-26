package load_balancer

import (
	"github.com/innotech/hydra/log"

	zmq "github.com/innotech/hydra/vendors/github.com/pebbe/zmq4"

	"errors"
	"runtime"
	"time"
)

type Requester interface {
	Close()
	Send([]byte, [][]byte) [][]byte
}

type Client struct {
	broker  string
	context *zmq.Context // Socket to broker
	poller  *zmq.Poller
	retries int // Request retries
	socket  *zmq.Socket
	timeout time.Duration // Request timeout
}

func NewClient(broker string, requestTimeout int) (cli *Client, err error) {
	cli = &Client{
		broker:  broker,
		timeout: time.Duration(requestTimeout) * time.Millisecond,
		// TODO: new
		retries: 3, //  Before we abandon
	}
	cli.context, err = zmq.NewContext()
	if err != nil {
		log.Fatal("LoadBalancer client creating context failed")
	}
	err = cli.ConnectToBroker()
	runtime.SetFinalizer(cli, (*Client).Close)
	return
}

// Connect or reconnect to broker.
func (c *Client) ConnectToBroker() (err error) {
	if c.socket != nil {
		// TODO: Maybe catch error
		c.socket.Close()
		c.socket = nil
	}

	c.socket, err = c.context.NewSocket(zmq.REQ)
	if err != nil {
		log.Fatal("LoadBalancer client ConnectToBroker() creating socket failed")
	}
	// TODO: Maybe set linger to 0
	// err = c.socket.SetLinger(0)

	c.poller = zmq.NewPoller()
	c.poller.Add(c.socket, zmq.POLLIN)

	log.Debugf("LoadBalancer client connecting to broker at %s...", c.broker)
	err = c.socket.Connect(c.broker)
	if err != nil {
		log.Fatal("LoadBalancer client ConnectToBroker() failed to connect to broker", c.broker)
	}

	return
}

func (c *Client) Close() (err error) {
	if c.socket != nil {
		err = c.socket.Close()
		c.socket = nil
	}
	if c.context != nil {
		err = c.context.Term()
		c.context = nil
	}
	return
}

//  Send sends a request to the broker and gets a
//  reply even if it has to retry several times. It returns the reply
//  message, or error if there was no reply after multiple attempts:
func (c *Client) Send(service []byte, request [][]byte) (reply [][]byte, err error) {
	req := append([][]byte{service}, request...)
	log.Debugf("LoadBalancer client send request to '%s' service: %q\n", service, req)
	for retries_left := c.retries; retries_left > 0; retries_left-- {
		_, err = c.socket.SendMessage(req)
		if err != nil {
			break
		}

		//  On any blocking call, libzmq will return -1 if there was an error
		var polled []zmq.Polled
		polled, err = c.poller.Poll(c.timeout)
		if err != nil {
			break //  Interrupted
		}

		if len(polled) > 0 {
			var msg [][]byte
			msg, err = c.socket.RecvMessageBytes(0)
			if err != nil {
				break
			}
			log.Debugf("LoadBalancer client received reply: %q\n", msg)
			if len(msg) < 1 {
				reply = [][]byte{}
				return
			}

			reply = msg
			return //  Success
		} else {
			log.Warn("LoadBalancer client no reply, reconnecting...")
			c.ConnectToBroker()
		}
	}
	if err == nil {
		err = errors.New("LoadBalancer client permanent error")
	}
	log.Debug("LoadBalancer client permanent error, abandoning")
	reply = [][]byte{}
	return
}
