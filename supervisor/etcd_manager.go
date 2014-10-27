package supervisor

import (
	"github.com/innotech/hydra/database/connector"
	"github.com/innotech/hydra/etcd"
	"github.com/innotech/hydra/log"
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"sync"
	"time"
)

type EtcdController interface {
	Restart(config *etcd_config.Config)
	Start(config *etcd_config.Config)
	Stop()
}

type EtcdManager struct {
	sync.RWMutex
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
	log.Debug("Starting etcd service")
	e.Lock()
	defer e.Unlock()
	etcd.EtcdFactory.Config(config)
	e.EtcdService = etcd.EtcdFactory.Build()
	go func() {
		e.EtcdService.Start()
	}()
	// TODO: Maybe change argument type to EtcdService interface
	connector.SetEtcdConnector(e.EtcdService)
	time.Sleep(time.Duration(200) * time.Millisecond)
}

func (e *EtcdManager) Stop() {
	log.Debug("Stopping etcd service")
	e.Lock()
	defer e.Unlock()
	e.EtcdService.Stop()
	time.Sleep(time.Duration(200) * time.Millisecond)
}
