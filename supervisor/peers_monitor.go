package supervisor

import (
	"github.com/innotech/hydra/log"

	"reflect"
	"time"
)

type FolderMonitor interface {
	Run(ch chan []Peer)
}

// const (
// 	ClusterKey                    string        = "cluster"
// 	AddrKey                       string        = "addr"
// 	PeerStateEnabled              string        = "enabled"
// 	DefaultRequestClusterInterval time.Duration = time.Duration(3) * time.Second
// )

type PeersMonitor struct {
	Peers                  []Peer
	etcdClient             EtcdRequester
	RequestClusterInterval time.Duration
}

func NewPeersMonitor(etcdClient EtcdRequester) *PeersMonitor {
	return &PeersMonitor{
		Peers:                  []Peer{},
		etcdClient:             etcdClient,
		RequestClusterInterval: DefaultRequestClusterInterval,
	}
}

func (p *PeersMonitor) Run(ch chan []Peer) {
	var res *Response
	var err error
	var newCLuster *PeerCluster
	for {
		res, err = p.etcdClient.Get(ClusterKey, true, true)
		if err == nil {
			newCLuster = NewPeerClusterFromNodes(res.Node.Nodes)
			if !reflect.DeepEqual(p.Peers, newCLuster.Peers) {
				p.Peers = newCLuster.Peers
				ch <- p.Peers
			}
		} else {
			log.Warn("Unreachable cluster container - thrown error:  " + err.Error())
		}
		time.Sleep(p.RequestClusterInterval)
	}
}
