package etcd

import (
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"
)

type EtcdBuilder interface {
	Build() *Etcd
	Config(config *etcd_config.Config)
}

type etcdFactory struct {
	etcdConfig *etcd_config.Config
}

var EtcdFactory *etcdFactory = new(etcdFactory)

func (e *etcdFactory) Config(config *etcd_config.Config) {
	e.etcdConfig = config
}

func (e *etcdFactory) Build() *Etcd {
	// TODO: refactor Etcd visibility
	etcd := New(e.etcdConfig)
	etcd.Load()
	return etcd
}
