package supervisor

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
		EtcdClient:   NewEtcdClient([]string{ownAddr}),
		ownAddr:      ownAddr,
		PeerCluster:  NewPeerCluster(peers),
		PeersMonitor: NewPeersMonitor(NewEtcdClient([]string{ownAddr}).WithMachineAddr(ownAddr)),
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
			c.PeerCluster.Peers = newPeers
		default:
			leader := c.searchForLeader()
			if leader != nil {
				foreignPeerCluster, err := c.getCluster(AddHttpProtocol((*leader).Addr))
				if err == nil {
					if c.shouldTryToJoinTheCluster((*leader).Addr, foreignPeerCluster) {
						// TODO: Get PeerAddr from foreignPeerCluster data
						ch <- AddHttpProtocol((*leader).PeerAddr)
						break OuterLoop
					}
				}
			}
			// TODO: Add sleep
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
	c.EtcdClient.WithMachineAddr(addr)
	res, err := c.EtcdClient.Get(ClusterKey, true, true)
	if err != nil {
		return nil, err
	}
	return NewPeerClusterFromNodes(res.Node.Nodes), nil
}

func (c *ClusterInspector) searchForLeader() *Peer {
	for c.PeerCluster.Reset(); c.PeerCluster.HasNext(); {
		peer, _ := c.PeerCluster.Next()
		addrURL := AddHttpProtocol(peer.Addr)
		peerLeader, err := c.getLeader(addrURL)
		if err == nil && peerLeader == AddHttpProtocol(peer.PeerAddr) {
			return &peer
		}
	}
	return nil
}

func (c *ClusterInspector) getLeader(addr string) (string, error) {
	c.EtcdClient.WithMachineAddr(addr)
	res, err := c.EtcdClient.BaseGet(LeaderKey)
	if err != nil {
		return "", err
	}
	leader := string(res.Body)
	return leader, nil
}
