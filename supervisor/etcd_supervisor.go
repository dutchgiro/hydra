package supervisor

import (
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"time"
)

type Supervisor interface {
	Run()
}

const (
	AddrKey                       string        = "addr"
	ClusterKey                    string        = "cluster"
	DefaultMaxAttemptsToSetState  uint          = 3
	DefaultRequestClusterInterval time.Duration = time.Duration(3) * time.Second
	LeaderKey                     string        = "leader"
	PeerAddrKey                   string        = "peerAddr"
	PeerStateEnabled              string        = "enabled"
	StateKey                      string        = "state"
)

type EtcdSupervisor struct {
	etcdConfig       *etcd_config.Config
	ClusterInspector ClusterAnalyzer
	EtcdManager      EtcdController
	StateManager     StateController
}

func NewEtcdSupervisor(config *etcd_config.Config) *EtcdSupervisor {
	configPeers := config.Peers
	config.Peers = []string{}
	return &EtcdSupervisor{
		etcdConfig: config,
		// ClusterInspector: NewClusterInspector(config.Addr, config.Peer.Addr, configPeers),
		ClusterInspector: NewClusterInspector(config.Addr, configPeers),
		EtcdManager:      NewEtcdManager(),
		StateManager: NewStateManager(
			time.Duration(config.Peer.HeartbeatTimeout*3)*time.Millisecond,
			NewEtcdClient([]string{config.Addr}).WithMachineAddr(config.Addr),
			config.Peer.Addr,
			// TODO: should be calculates from first argument
			uint64(config.Peer.HeartbeatTimeout*6),
		),
	}
}

func (e *EtcdSupervisor) Run() {
	e.EtcdManager.Start(e.etcdConfig)

	stateManagerChannel := make(chan StateControllerState)
	e.StateManager.Run(stateManagerChannel)
	clusterInspectorChannel := make(chan string)
	e.ClusterInspector.Run(clusterInspectorChannel)

	var newLeader string
	for {
		select {
		case <-stateManagerChannel:
			e.etcdConfig.Peers = []string{}
			e.EtcdManager.Restart(e.etcdConfig)
		case newLeader = <-clusterInspectorChannel:
			e.etcdConfig.Peers = []string{newLeader}
			e.EtcdManager.Restart(e.etcdConfig)
		}
	}
}
