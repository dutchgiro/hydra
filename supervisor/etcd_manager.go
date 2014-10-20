package supervisor

import (
	"github.com/innotech/hydra/database/connector"
	"github.com/innotech/hydra/etcd"
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"
)

type EtcdController interface {
	Restart(config *etcd_config.Config)
	Start(config *etcd_config.Config)
	Stop()
}

type EtcdManager struct {
	EtcdService etcd.EtcdService
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
	// TODO: Maybe change argument type to EtcdService interface
	connector.SetEtcdConnector(e.EtcdService)
}

func (e *EtcdManager) Stop() {
	e.EtcdService.Stop()
}
