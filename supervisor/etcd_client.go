package supervisor

import (
	"net"
	"net/http"
	"time"
)

type Node struct {
	Key           string     `json:"key, omitempty"`
	Value         string     `json:"value,omitempty"`
	Dir           bool       `json:"dir,omitempty"`
	Expiration    *time.Time `json:"expiration,omitempty"`
	TTL           int64      `json:"ttl,omitempty"`
	Nodes         Nodes      `json:"nodes,omitempty"`
	ModifiedIndex uint64     `json:"modifiedIndex,omitempty"`
	CreatedIndex  uint64     `json:"createdIndex,omitempty"`
}

type Nodes []*Node

type Response struct {
	Action    string `json:"action"`
	Node      *Node  `json:"node"`
	PrevNode  *Node  `json:"prevNode,omitempty"`
	EtcdIndex uint64 `json:"etcdIndex"`
	RaftIndex uint64 `json:"raftIndex"`
	RaftTerm  uint64 `json:"raftTerm"`
}

type EtcdRequester interface {
	CompareAndSwap()
}

type Client struct {
	config      Config   `json:"config"`
	cluster     *Cluster `json:"cluster"`
	httpClient  *http.Client
	persistence io.Writer
	cURLch      chan string
	// CheckRetry can be used to control the policy for failed requests
	// and modify the cluster if needed.
	// The client calls it before sending requests again, and
	// stops retrying if CheckRetry returns some error. The cases that
	// this function needs to handle include no response and unexpected
	// http status code of response.
	// If CheckRetry is nil, client will call the default one
	// `DefaultCheckRetry`.
	// Argument cluster is the etcd.Cluster object that these requests have been made on.
	// Argument numReqs is the number of http.Requests that have been made so far.
	// Argument lastResp is the http.Responses from the last request.
	// Argument err is the reason of the failure.
	CheckRetry func(cluster *Cluster, numReqs int,
		lastResp http.Response, err error) error
}

// NewClient create a basic client that is configured to be used
// with the given machine list.
func NewClient(machines []string) *Client {
	config := Config{
		// default timeout is one second
		DialTimeout: time.Second,
		// default consistency level is STRONG
		Consistency: STRONG_CONSISTENCY,
	}

	client := &Client{
		cluster: NewCluster(machines),
		config:  config,
	}

	client.initHTTPClient()
	client.saveConfig()

	return client
}

// initHTTPClient initializes a HTTP client for etcd client
func (c *Client) initHTTPClient() {
	tr := &http.Transport{
		Dial: c.dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	c.httpClient = &http.Client{Transport: tr}
}

// dial attempts to open a TCP connection to the provided address, explicitly
// enabling keep-alives with a one-second interval.
func (c *Client) dial(network, addr string) (net.Conn, error) {
	conn, err := net.DialTimeout(network, addr, c.config.DialTimeout)
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

// SendRequest sends a HTTP request and returns a Response as defined by etcd
func (c *Client) SendRequest(rr *RawRequest) (*RawResponse, error) {

	var req *http.Request
	var resp *http.Response
	var httpPath string
	var err error
	var respBody []byte

	var numReqs = 1

	checkRetry := c.CheckRetry
	if checkRetry == nil {
		checkRetry = DefaultCheckRetry
	}

	cancelled := make(chan bool, 1)
	reqLock := new(sync.Mutex)

	if rr.Cancel != nil {
		cancelRoutine := make(chan bool)
		defer close(cancelRoutine)

		go func() {
			select {
			case <-rr.Cancel:
				cancelled <- true
				logger.Debug("send.request is cancelled")
			case <-cancelRoutine:
				return
			}

			// Repeat canceling request until this thread is stopped
			// because we have no idea about whether it succeeds.
			for {
				reqLock.Lock()
				c.httpClient.Transport.(*http.Transport).CancelRequest(req)
				reqLock.Unlock()

				select {
				case <-time.After(100 * time.Millisecond):
				case <-cancelRoutine:
					return
				}
			}
		}()
	}

	// If we connect to a follower and consistency is required, retry until
	// we connect to a leader
	sleep := 25 * time.Millisecond
	maxSleep := time.Second
	for attempt := 0; ; attempt++ {
		if attempt > 0 {
			select {
			case <-cancelled:
				return nil, ErrRequestCancelled
			case <-time.After(sleep):
				sleep = sleep * 2
				if sleep > maxSleep {
					sleep = maxSleep
				}
			}
		}

		logger.Debug("Connecting to etcd: attempt", attempt+1, "for", rr.RelativePath)

		if rr.Method == "GET" && c.config.Consistency == WEAK_CONSISTENCY {
			// If it's a GET and consistency level is set to WEAK,
			// then use a random machine.
			httpPath = c.getHttpPath(true, rr.RelativePath)
		} else {
			// Else use the leader.
			httpPath = c.getHttpPath(false, rr.RelativePath)
		}

		// Return a cURL command if curlChan is set
		if c.cURLch != nil {
			command := fmt.Sprintf("curl -X %s %s", rr.Method, httpPath)
			for key, value := range rr.Values {
				command += fmt.Sprintf(" -d %s=%s", key, value[0])
			}
			c.sendCURL(command)
		}

		logger.Debug("send.request.to ", httpPath, " | method ", rr.Method)

		reqLock.Lock()
		if rr.Values == nil {
			if req, err = http.NewRequest(rr.Method, httpPath, nil); err != nil {
				return nil, err
			}
		} else {
			body := strings.NewReader(rr.Values.Encode())
			if req, err = http.NewRequest(rr.Method, httpPath, body); err != nil {
				return nil, err
			}

			req.Header.Set("Content-Type",
				"application/x-www-form-urlencoded; param=value")
		}
		reqLock.Unlock()

		resp, err = c.httpClient.Do(req)
		defer func() {
			if resp != nil {
				resp.Body.Close()
			}
		}()

		// If the request was cancelled, return ErrRequestCancelled directly
		select {
		case <-cancelled:
			return nil, ErrRequestCancelled
		default:
		}

		numReqs++

		// network error, change a machine!
		if err != nil {
			logger.Debug("network error:", err.Error())
			lastResp := http.Response{}
			if checkErr := checkRetry(c.cluster, numReqs, lastResp, err); checkErr != nil {
				return nil, checkErr
			}

			c.cluster.switchLeader(attempt % len(c.cluster.Machines))
			continue
		}

		// if there is no error, it should receive response
		logger.Debug("recv.response.from", httpPath)

		if validHttpStatusCode[resp.StatusCode] {
			// try to read byte code and break the loop
			respBody, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				logger.Debug("recv.success.", httpPath)
				break
			}
			// ReadAll error may be caused due to cancel request
			select {
			case <-cancelled:
				return nil, ErrRequestCancelled
			default:
			}
		}

		// if resp is TemporaryRedirect, set the new leader and retry
		if resp.StatusCode == http.StatusTemporaryRedirect {
			u, err := resp.Location()

			if err != nil {
				logger.Warning(err)
			} else {
				// Update cluster leader based on redirect location
				// because it should point to the leader address
				c.cluster.updateLeaderFromURL(u)
				logger.Debug("recv.response.relocate", u.String())
			}
			resp.Body.Close()
			continue
		}

		if checkErr := checkRetry(c.cluster, numReqs, *resp,
			errors.New("Unexpected HTTP status code")); checkErr != nil {
			return nil, checkErr
		}
		resp.Body.Close()
	}

	r := &RawResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Header:     resp.Header,
	}

	return r, nil
}

// put issues a PUT request
func (c *Client) put(key string, value string, ttl uint64,
	options Options) (*RawResponse, error) {

	logger.Debugf("put %s, %s, ttl: %d, [%s]", key, value, ttl, c.cluster.Leader)
	p := keyToPath(key)

	str, err := options.toParameters(VALID_PUT_OPTIONS)
	if err != nil {
		return nil, err
	}
	p += str

	req := NewRawRequest("PUT", p, buildValues(value, ttl), nil)
	resp, err := c.SendRequest(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) CompareAndSwap(key string, value string, ttl uint64,
	prevValue string, prevIndex uint64) (*Response, error) {
	raw, err := c.RawCompareAndSwap(key, value, ttl, prevValue, prevIndex)
	if err != nil {
		return nil, err
	}

	return raw.Unmarshal()
}

func (c *Client) RawCompareAndSwap(key string, value string, ttl uint64,
	prevValue string, prevIndex uint64) (*RawResponse, error) {
	if prevValue == "" && prevIndex == 0 {
		return nil, fmt.Errorf("You must give either prevValue or prevIndex.")
	}

	options := Options{}
	if prevValue != "" {
		options["prevValue"] = prevValue
	}
	if prevIndex != 0 {
		options["prevIndex"] = prevIndex
	}

	raw, err := c.put(key, value, ttl, options)

	if err != nil {
		return nil, err
	}

	return raw, err
}
