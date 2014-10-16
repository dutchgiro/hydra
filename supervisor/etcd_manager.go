package supervisor

import (
	. "github.com/innotech/hydra/etcd"
	"github.com/innotech/hydra/log"

	"os"
)

type EtcdController interface {
	Restart()
	Start()
	Stop()
}

type EtcdManager struct {
	etcdService *etcd.Etcd
}

func NewEtcdManager() {
	return new(EtcdManager)
}

func (e *EtcdManager) Restart() {
	e.Stop()
	e.Start()
}

func (e *EtcdManager) Start() {
	EtcdFactory.Config(config)
	e.etcdService = EtcdFactory.Build()
	e.etcdService.Start()
}

func (e *EtcdManager) Stop() {
	e.etcdService.Stop()
	e.removeEtcdDataFiles()
}

func (e *EtcdManager) removeEtcdDataFiles() {
	if err := os.RemoveAll(e.etcdService.Config.DataDir); err != nil {
		log.Fatal(err.Error())
	}
}
