package supervisor

type ClusterAnalyzer interface {
	Run(ch chan string)
}

const (
	LeaderKey   string = "leader"
	MachinesKey string = "machines"
)

type ClusterInspector struct {
	EtcdClient   EtcdRequester
	configPeers  []string
	myPeerAddr   string
	PeerCluster  *PeerCluster
	PeersMonitor FolderMonitor
}

func NewClusterInspector(ownAddr, ownPeerAddr string, knownPeers []string) *ClusterInspector {
	return &ClusterInspector{
		configPeers:  []string{},
		EtcdClient:   NewEtcdClient([]string{ownAddr}),
		myPeerAddr:   ownPeerAddr,
		PeerCluster:  NewPeerCluster([]Peer{}),
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
					if c.shouldTryToJoinTheCluster((*leader).PeerAddr, foreignPeerCluster) {
						ch <- AddHttpProtocol((*leader).PeerAddr)
						break OuterLoop
					}
				}
			}
			// TODO: Add sleep
		}
	}
}

func (c *ClusterInspector) shouldTryToJoinTheCluster(foreignLeaderPeerAddr string, foreignCluster *PeerCluster) bool {
	if len(foreignCluster.GetEnabledPeers()) > len(c.PeerCluster.GetEnabledPeers()) {
		return true
	} else if len(foreignCluster.GetEnabledPeers()) < len(c.PeerCluster.GetEnabledPeers()) {
		return false
	} else if c.iHavePriorityOverCluster(foreignLeaderPeerAddr, foreignCluster) {
		return false
	}
	return true
}

func (c *ClusterInspector) iHavePriorityOverCluster(foreignLeaderPeerAddr string, foreignCluster *PeerCluster) bool {
	myPositionInForeignCluster, err := foreignCluster.GetPeerPosition(c.myPeerAddr)
	if err != nil {
		return false
	}
	foreignLeaderPositionInOwnCluster, err := c.PeerCluster.GetPeerPosition(foreignLeaderPeerAddr)
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
