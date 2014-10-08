package supervisor

import (
	"github.com/innotech/hydra/etcd"
)

type EtcdController interface {
	Restart()
	Start()
	Stop()
}

type EtcdManager struct {
	etcdService *etcd.Etcd
	peers       []string
}

func (e *EtcdManager) Restart() {
	e.Stop()
	e.Start()
}

func (e *EtcdManager) SavePeers(clusterPeers []string) {
	e.peers = clusterPeers
}

func (e *EtcdManager) Start() {
	// TODO: drop withEtcdServer argument
	e.etcdService.Start()
}

func (e *EtcdManager) Stop() {
	e.etcdService.Stop()
	// var err error
	// err = etcdService.EtcdServerListener.Close()
	// if err != nil {
	// 	log.Println(errClose.Error())
	// }
	// err = etcdService.PeerServerListener.Close()
	// if err != nil {
	// 	log.Println(errClose.Error())
	// }
}
