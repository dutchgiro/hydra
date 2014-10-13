package supervisor

type ClusterAnalyzer interface {
	Run()
}

const (
	// TODO: unify etcd constants
	LeaderKey string = "leader"
)

type ClusterInspector struct {
	etcdClient   EtcdRequester
	configPeers  []string
	peers        PeerCluster
	PeersMonitor FolderMonitor
}

// TODO: etcdClient not necessary
func NewClusterInspector(selfAddr, selfPeerAddr string, knownPeers []string) *ClusterInspector {
	return &ClusterInspector{
		etcdClient:   NewEtcdClient([]string{selfAddr}),
		configPeers:  []string{},
		PeersMonitor: NewPeersMonitor(NewEtcdClient([]string{selfAddr}).WithMachineAddr(selfAddr)),
	}
}

func (c *ClusterInspector) Run() {
	// var peerCluster []Peer
	peersMonitorChannel := make(chan []Peer)
	go c.PeersMonitor.Run(peersMonitorChannel)
	// OuterLoop:
	for {
		select {
		// case <-e.executionChannel:
		// 	fmt.Println("------->>> BREAK LOOP")
		// 	break OuterLoop
		// case peerCluster = <-peersMonitorChannel:
		case _ = <-peersMonitorChannel:

			break
			// default:
			// 	leader := c.searchForLeader()
			// 	time.Sleep(e.durationBetweenPublicationsState)
		}
	}
}

func (c *ClusterInspector) mergeCluster() {

}

func (c *ClusterInspector) searchForLeader() {
	// for peer := range peers {
	// 	peerLeader, err := c.getPeerLeader(peer)
	// 	if err == nil && peerLeader != "" {

	// 	}
	// }
}

func (c *ClusterInspector) getPeerLeader(peerAddr string) (string, error) {
	res, err := c.etcdClient.WithMachineAddr(peerAddr).BaseGet(LeaderKey)
	if err != nil {
		return "", err
	}
	leader := string(res.Body)
	return leader, nil
}

func (c *ClusterInspector) refreshCluster() {
	// res, err := c.etcdClient.WithMachineAddr(peerAddr).BaseGet(LeaderKey)
}
