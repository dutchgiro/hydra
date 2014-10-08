package supervisor

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type EtcdRequester interface {
	CompareAndSwap(key string, value string, ttl uint64,
		prevValue string, prevIndex uint64) (*Response, error)
	Set(key string, value string, ttl uint64) (*Response, error)
}

// See SetConsistency for how to use these constants.
const (
	version = "v2"

	// Using strings rather than iota because the consistency level
	// could be persisted to disk, so it'd be better to use
	// human-readable values.
	STRONG_CONSISTENCY = "STRONG"
	WEAK_CONSISTENCY   = "WEAK"
)

const (
	defaultBufferSize = 10
)

type Config struct {
	CertFile    string        `json:"certFile"`
	KeyFile     string        `json:"keyFile"`
	CaCertFile  []string      `json:"caCertFiles"`
	DialTimeout time.Duration `json:"timeout"`
	Consistency string        `json:"consistency"`
}

type EtcdClient struct {
	machineAddr string

	config      Config   `json:"config"`
	cluster     *Cluster `json:"cluster"`
	httpClient  *http.Client
	persistence io.Writer
	cURLch      chan string
	// 	// CheckRetry can be used to control the policy for failed requests
	// 	// and modify the cluster if needed.
	// 	// The client calls it before sending requests again, and
	// 	// stops retrying if CheckRetry returns some error. The cases that
	// 	// this function needs to handle include no response and unexpected
	// 	// http status code of response.
	// 	// If CheckRetry is nil, client will call the default one
	// 	// `DefaultCheckRetry`.
	// 	// Argument cluster is the etcd.Cluster object that these requests have been made on.
	// 	// Argument numReqs is the number of http.Requests that have been made so far.
	// 	// Argument lastResp is the http.Responses from the last request.
	// 	// Argument err is the reason of the failure.
	CheckRetry func(cluster *Cluster, numReqs int,
		lastResp http.Response, err error) error
}

// NewClient create a basic client that is configured to be used
// with the given machine list.
func NewEtcdClient(machines []string) *EtcdClient {
	config := Config{
		// default timeout is one second
		DialTimeout: time.Second,
		// default consistency level is STRONG
		Consistency: STRONG_CONSISTENCY,
	}

	client := &EtcdClient{
		cluster: NewCluster(machines),
		config:  config,
	}

	client.initHTTPClient()
	client.saveConfig()

	return client
}

// initHTTPClient initializes a HTTP client for etcd client
func (e *EtcdClient) initHTTPClient() {
	tr := &http.Transport{
		Dial: e.dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	e.httpClient = &http.Client{Transport: tr}
}

// dial attempts to open a TCP connection to the provided address, explicitly
// enabling keep-alives with a one-second interval.
func (e *EtcdClient) dial(network, addr string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, addr, e.config.DialTimeout)
	if err != nil {
		return nil, err
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, errors.New("Failed type-assertion of net.Conn as *net.TCPConn")
	}

	// Keep TCP alive to check whether or not the remote machine is down
	if err = tcpConn.SetKeepAlive(true); err != nil {
		return nil, err
	}

	if err = tcpConn.SetKeepAlivePeriod(time.Second); err != nil {
		return nil, err
	}

	return tcpConn, nil
}

func (e *EtcdClient) sendCURL(command string) {
	go func() {
		select {
		case e.cURLch <- command:
		default:
		}
	}()
}

// saveConfig saves the current config using c.persistence.
func (e *EtcdClient) saveConfig() error {
	if e.persistence != nil {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}

		_, err = e.persistence.Write(b)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *EtcdClient) WithMachineAddr(machineAddr string) *EtcdClient {
	e.machineAddr = machineAddr

	return e
}

// Compare And Swap

type KeyExistence int

const (
	False = iota
	True
	Unknow
)

func (e *EtcdClient) CompareAndSwap(key string, value string, ttl uint64,
	prevValue string, prevIndex uint64, prevExist KeyExistence) (*Response, error) {
	raw, err := e.RawCompareAndSwap(key, value, ttl, prevValue, prevIndex, prevExist)
	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

func (e *EtcdClient) RawCompareAndSwap(key string, value string, ttl uint64,
	prevValue string, prevIndex uint64, prevExist KeyExistence) (*RawResponse, error) {
	if prevValue == "" && prevIndex == 0 && prevExist == Unknow {
		return nil, fmt.Errorf("You must give either prevValue or prevIndex.")
	}

	options := Options{}
	if prevValue != "" {
		options["prevValue"] = prevValue
	}
	if prevIndex != 0 {
		options["prevIndex"] = prevIndex
	}
	if prevExist != Unknow {
		var exist bool
		if prevExist == False {
			exist = false
		} else if prevExist == True {
			exist = true
		}
		options["prevExist"] = exist
	}

	raw, err := e.put(key, value, ttl, options)

	if err != nil {
		return nil, err
	}

	return raw, err
}
