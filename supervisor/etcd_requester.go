package supervisor

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"
// )

// type Node struct {
// 	Key           string     `json:"key, omitempty"`
// 	Value         string     `json:"value,omitempty"`
// 	Dir           bool       `json:"dir,omitempty"`
// 	Expiration    *time.Time `json:"expiration,omitempty"`
// 	TTL           int64      `json:"ttl,omitempty"`
// 	Nodes         Nodes      `json:"nodes,omitempty"`
// 	ModifiedIndex uint64     `json:"modifiedIndex,omitempty"`
// 	CreatedIndex  uint64     `json:"createdIndex,omitempty"`
// }

// type Nodes []*Node

// type Response struct {
// 	Action    string `json:"action"`
// 	Node      *Node  `json:"node"`
// 	PrevNode  *Node  `json:"prevNode,omitempty"`
// 	EtcdIndex uint64 `json:"etcdIndex"`
// 	RaftIndex uint64 `json:"raftIndex"`
// 	RaftTerm  uint64 `json:"raftTerm"`
// }

// // TODO: Change name to Monitor
// type EtcdRequester interface {
// 	Put(key string, value string) error
// }

// type SupervisorRequester struct {
// 	etcdAddr string
// }

// func NewSupervisorRequester(etcdAddr string) *SupervisorRequester {
// 	return &SupervisorRequester{
// 		etcdAddr: etcdAddr,
// 	}
// }

// func (s *SupervisorRequester) SetNodeState(hydraClusterRootPath, hydraNodeId, state string) (Response, error) {
// 	client := &http.Client{}
// 	// TODO: Replace for state (master, master', slave)
// 	reqBody := "value=" + state
// 	req, err := http.NewRequest("PUT", hydraClusterRootPath+hydraNodeId, body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	resBody, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var etcdResponse Response
// 	err = json.Unmarshal(resBody, &etcdResponse)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return etcdResponse, nil
// }
