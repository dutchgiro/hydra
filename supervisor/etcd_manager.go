package supervisor

import (
	"github.com/innotech/hydra/etcd"
	"github.com/innotech/hydra/log"
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"os"
)

type EtcdController interface {
	Restart()
	Start()
	Stop()
}

type EtcdManager struct {
	EtcdService *etcd.Etcd
}

func NewEtcdManager() *EtcdManager {
	return new(EtcdManager)
}

func (e *EtcdManager) Restart(config *etcd_config.Config) {
	e.Stop()
	e.Start(config)
}

func (e *EtcdManager) Start(config *etcd_config.Config) {
	etcd.EtcdFactory.Config(config)
	e.EtcdService = etcd.EtcdFactory.Build()
	e.EtcdService.Start()
}

func (e *EtcdManager) Stop() {
	e.EtcdService.Stop()
	e.removeEtcdDataFiles()
}

func (e *EtcdManager) removeEtcdDataFiles() {
	if err := os.RemoveAll(e.EtcdService.Config.DataDir); err != nil {
		log.Fatal(err.Error())
	}
}
