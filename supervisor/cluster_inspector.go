package supervisor

type ClusterAnalyzer interface {
	Run()
}

const (
	// TODO: unify etcd constants
	LeaderKey string = "leader"
)

type ClusterInspector struct {
	EtcdClient   EtcdRequester
	configPeers  []string
	peerCluster  PeerCluster
	PeersMonitor FolderMonitor
}

// TODO: etcdClient not necessary
func NewClusterInspector(selfAddr, selfPeerAddr string, knownPeers []string) *ClusterInspector {
	return &ClusterInspector{
		EtcdClient:   NewEtcdClient([]string{selfAddr}),
		configPeers:  []string{},
		PeersMonitor: NewPeersMonitor(NewEtcdClient([]string{selfAddr}).WithMachineAddr(selfAddr)),
	}
}

func (c *ClusterInspector) Run() {
	var newPeers []Peer
	peersMonitorChannel := make(chan []Peer)
	go c.PeersMonitor.Run(peersMonitorChannel)
	// OuterLoop:
	for {
		select {
		// case <-e.executionChannel:
		// 	fmt.Println("------->>> BREAK LOOP")
		// 	break OuterLoop
		// case peerCluster = <-peersMonitorChannel:
		case newPeers = <-peersMonitorChannel:
			c.peerCluster.SetPeers(newPeers)
		default:
			// leader := c.searchForLeader()
			_ = c.searchForLeader()
			// time.Sleep(e.durationBetweenPublicationsState)
		}
	}
}

func (c *ClusterInspector) mergeCluster() {

}

func (c *ClusterInspector) searchForLeader() string {
	// TODO: Reset iterator
	for /* reset iterator*/ c.peerCluster.HasNext() {
		peer, _ := c.peerCluster.Next()
		peerLeader, err := c.getPeerLeader(peer.Addr)
		if err == nil && peerLeader != "" {

		}
	}
	return ""
}

func (c *ClusterInspector) getPeerLeader(addr string) (string, error) {
	res, err := c.EtcdClient.WithMachineAddr(addr).BaseGet(LeaderKey)
	if err != nil {
		return "", err
	}
	leader := string(res.Body)
	return leader, nil
}

func (c *ClusterInspector) refreshCluster() {
	// res, err := c.etcdClient.WithMachineAddr(peerAddr).BaseGet(LeaderKey)
}
