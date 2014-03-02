package connector

import (
	"time"

	"github.com/innotech/hydra/etcd"

	// "github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/coreos/raft"
)

type EtcdDriver interface {
	Delete(key string, dir, recursive bool) error
	Get(key string, recursive bool, sort bool) []interface{}
	Set(key string, dir bool, value string, expireTime time.Time) error
}

type EtcdConnector struct {
	etcd *etcd.Etcd
}

var e *EtcdConnector

func SetEtcdConnector(etcd *etcd.Etcd) {
	e = new(EtcdConnector)
	e.etcd = etcd
}

func GetEtcdConnector() *EtcdConnector {
	return e
}

func (e EtcdConnector) Delete(key string, dir, recursive bool) error {
	return nil

}

func (e EtcdConnector) Get(key string, recursive bool, sort bool) []interface{} {
	return nil
}

func (e EtcdConnector) Set(key string, dir bool, value string, expireTime time.Time) error {
	c := e.etcd.EtcdServer.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
	// c := e.etcd.Store().CommandFactory().CreateSetCommand(key, dir, value, expireTime)
	return e.dispatch(c)
}

func (e EtcdConnector) dispatch(c raft.Command) error {
	ps := e.etcd.PeerServer
	if ps.RaftServer().State() == raft.Leader {
		result, err := ps.RaftServer().Do(c)
		if err != nil {
			return err
		}

		if result == nil {
			return nil
		}

		return nil

	} else {
		leader := ps.RaftServer().Leader()

		if leader == "" {
			return nil
		}

		return nil
	}
}