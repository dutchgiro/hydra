package supervisor

import (
	"github.com/innotech/hydra/log"

	"reflect"
	"time"
)

type FolderMonitor interface {
	Run(ch chan PeerCluster)
}

const (
	ClusterKey string = "cluster"
	// TODO: unify every cluster management constant
	// StateKey                      string        = "state"
	PeerStateEnabled              string        = "enabled"
	DefaultRequestClusterInterval time.Duration = time.Duration(3) * time.Second
)

type Peer struct {
	Id    string
	State string
}

type PeerCluster []Peer

type PeersMonitor struct {
	Cluster                PeerCluster
	etcdClient             EtcdRequester
	RequestClusterInterval time.Duration
}

func NewPeersMonitor(etcdClient EtcdRequester) *PeersMonitor {
	return &PeersMonitor{
		Cluster:                []Peer{},
		etcdClient:             etcdClient,
		RequestClusterInterval: DefaultRequestClusterInterval,
	}
}

func (p *PeersMonitor) Run(ch chan PeerCluster) {
	var res *Response
	var err error
	for {
		res, err = p.etcdClient.Get(ClusterKey, true, true)
		if err == nil {
			resPeerCluster := p.parseRawPeerCluster(res)
			if !reflect.DeepEqual(p.Cluster, resPeerCluster) {
				p.Cluster = resPeerCluster
				ch <- p.Cluster
			}
		} else {
			log.Warn("Unreachable cluster container - thrown error:  " + err.Error())
		}
		time.Sleep(p.RequestClusterInterval)
	}
}

func (p *PeersMonitor) parseRawPeerCluster(rawPeerCluster *Response) PeerCluster {
	var peerCluster PeerCluster = []Peer{}
	rawPeers := rawPeerCluster.Node.Nodes
	for _, rawPeer := range rawPeers {
		peerCluster = append(peerCluster, p.parseRawPeer(rawPeer))
	}
	return peerCluster
}

func (p *PeersMonitor) parseRawPeer(rawPeer *Node) Peer {
	peer := Peer{Id: rawPeer.Key}
	for _, attr := range rawPeer.Nodes {
		switch attr.Key {
		case StateKey:
			peer.State = attr.Value
		}
	}
	return peer
}
