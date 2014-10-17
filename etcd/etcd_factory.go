package etcd

import (
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"
)

type EtcdBuilder interface {
	Build() EtcdService
	Config(config *etcd_config.Config)
}

type etcdFactory struct {
	etcdConfig *etcd_config.Config
}

var EtcdFactory EtcdBuilder = new(etcdFactory)

func (e *etcdFactory) Config(config *etcd_config.Config) {
	e.etcdConfig = config
}

func (e *etcdFactory) Build() EtcdService {
	// TODO: refactor Etcd visibility
	etcd := New(e.etcdConfig)
	etcd.Load()
	return etcd
}
