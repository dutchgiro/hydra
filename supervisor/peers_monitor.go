package supervisor

import (
	"time"
)

type FolderMonitor interface {
	Run()
}

const (
	DefaultRequestClusterInterval time.Duration = time.Duration(3) * time.Second
)

type PeersMonitor struct {
	etcdClient             EtcdRequester
	RequestClusterInterval time.Duration
}

func NewPeersMonitor(etcdClient EtcdRequester) *PeersMonitor {
	return &PeersMonitor{
		etcdClient:             etcdClient,
		RequestClusterInterval: DefaultRequestClusterInterval,
	}
}

func (p *PeersMonitor) Run() {
	for {
		// res, err := p.etcdClient.Get("cluster", true, true)
		_, _ = p.etcdClient.Get("cluster", true, true)
		// if err == nil {
		// 	p.compareCluster(res.Node)
		// } else {
		// 	// TODO: log
		// }
		time.Sleep(p.RequestClusterInterval)
	}
}

func (p *PeersMonitor) compareCluster(cluster *Node) {

}
