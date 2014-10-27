package supervisor

import (
	"github.com/innotech/hydra/log"

	"time"
)

type ClusterAnalyzer interface {
	Run(ch chan string)
}

type ClusterInspector struct {
	EtcdClient   EtcdRequester
	configPeers  []string
	ownAddr      string
	PeerCluster  *PeerCluster
	PeersMonitor FolderMonitor
}

// func NewClusterInspector(ownAddr, ownPeerAddr string, knownPeers []string) *ClusterInspector {
func NewClusterInspector(ownAddr string, knownPeers []string) *ClusterInspector {
	peers := []Peer{}
	for i := 0; i < len(knownPeers); i++ {
		peers = append(peers, Peer{Addr: knownPeers[i]})
	}
	return &ClusterInspector{
		configPeers:  []string{},
		EtcdClient:   NewEtcdClient([]string{AddHttpProtocol(ownAddr)}),
		ownAddr:      ownAddr,
		PeerCluster:  NewPeerCluster(peers),
		PeersMonitor: NewPeersMonitor(NewEtcdClient([]string{AddHttpProtocol(ownAddr)}).WithMachineAddr(AddHttpProtocol(ownAddr))),
	}
}

func (c *ClusterInspector) Run(ch chan string) {
	var newPeers []Peer
	peersMonitorChannel := make(chan []Peer)
	go c.PeersMonitor.Run(peersMonitorChannel)
OuterLoop:
	for {
		select {
		case newPeers = <-peersMonitorChannel:
			log.Debug("Peer monitor found a change in cluster")
			c.PeerCluster.Peers = newPeers
		default:
			log.Debug("Search for leader")
			leader := c.searchForLeader()
			if leader != nil {
				foreignPeerCluster, err := c.getCluster(AddHttpProtocol((*leader).Addr))
				if err != nil {
					break
				}
				if c.shouldTryToJoinTheCluster((*leader).Addr, foreignPeerCluster) {
					ch <- AddHttpProtocol((*leader).PeerAddr)
					break OuterLoop
				}
			}
			// TODO: configure sleep
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}

func (c *ClusterInspector) shouldTryToJoinTheCluster(foreignLeaderAddr string, foreignCluster *PeerCluster) bool {
	if len(foreignCluster.GetEnabledPeers()) > len(c.PeerCluster.GetEnabledPeers()) {
		return true
	} else if len(foreignCluster.GetEnabledPeers()) < len(c.PeerCluster.GetEnabledPeers()) {
		return false
	} else if c.iHavePriorityOverCluster(foreignLeaderAddr, foreignCluster) {
		return false
	}
	return true
}

func (c *ClusterInspector) iHavePriorityOverCluster(foreignLeaderAddr string, foreignCluster *PeerCluster) bool {
	myPositionInForeignCluster, err := foreignCluster.GetPeerPosition(c.ownAddr)
	if err != nil {
		return false
	}
	foreignLeaderPositionInOwnCluster, err := c.PeerCluster.GetPeerPosition(foreignLeaderAddr)
	if err != nil {
		return true
	}
	if myPositionInForeignCluster <= foreignLeaderPositionInOwnCluster {
		return true
	}
	return false
}

func (c *ClusterInspector) getCluster(addr string) (*PeerCluster, error) {
	log.Debug("Call to getCluster with address: " + addr)
	c.EtcdClient.WithMachineAddr(addr)
	res, err := c.EtcdClient.Get(ClusterKey, true, true)
	if err != nil {
		log.Debug("Unable to get Cluster from peer: " + addr)
		return nil, err
	}
	return NewPeerClusterFromNodes(res.Node.Nodes), nil
}

func (c *ClusterInspector) searchForLeader() *Peer {
	for c.PeerCluster.Reset(); c.PeerCluster.HasNext(); {
		peer, _ := c.PeerCluster.Next()
		// TODO: Test
		if peer.Addr == c.ownAddr {
			continue
		}
		addrURL := AddHttpProtocol(peer.Addr)
		if peer.PeerAddr == "" {
			peerCluster, err := c.getCluster(addrURL)
			if err != nil {
				continue
			}
			p, err := peerCluster.GetPeerByAddr(peer.Addr)
			if err != nil || p.PeerAddr == "" {
				log.Debugf("Invalid peer in cluster: %#v", peerCluster)
				continue
			}
			peer.PeerAddr = p.PeerAddr
		}
		peerLeader, err := c.getLeader(addrURL)
		if err == nil && peerLeader == AddHttpProtocol(peer.PeerAddr) {
			return &peer
		}
	}
	return nil
}

func (c *ClusterInspector) getLeader(addr string) (string, error) {
	log.Debug("Call to getLeader with address: " + addr)
	log.Info("Call to getLeader with address: " + addr)
	c.EtcdClient.WithMachineAddr(addr)
	res, err := c.EtcdClient.BaseGet(LeaderKey)
	if err != nil {
		log.Debug("Unable to get leader from peer: " + addr)
		return "", err
	}
	leader := string(res.Body)
	return leader, nil
}
