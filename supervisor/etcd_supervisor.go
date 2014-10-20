package supervisor

import (
	hydra_config "github.com/innotech/hydra/config"

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
	PeerStateEnabled              string        = "enabled"
	StateKey                      string        = "state"
)

type EtcdSupervisor struct {
	hydraConfig      *hydra_config.Config
	ClusterInspector ClusterAnalyzer
	EtcdManager      EtcdController
	StateManager     StateController
}

func NewEtcdSupervisor(config *hydra_config.Config) *EtcdSupervisor {
	return &EtcdSupervisor{
		hydraConfig:      config,
		ClusterInspector: NewClusterInspector(config.EtcdAddr, config.Peer.Addr, config.Peers),
		EtcdManager:      NewEtcdManager(),
		StateManager: NewStateManager(
			time.Duration(config.Peer.HeartbeatTimeout*3)*time.Millisecond,
			NewEtcdClient([]string{config.EtcdAddr}).WithMachineAddr(config.EtcdAddr),
			config.Peer.Addr,
			// TODO: should be calculates from first argument
			uint64(config.Peer.HeartbeatTimeout*6),
		),
	}
}

func (e *EtcdSupervisor) Run() {
	stateManagerChannel := make(chan StateControllerState)
	e.StateManager.Run(stateManagerChannel)
	clusterInspectorChannel := make(chan string)
	e.ClusterInspector.Run(clusterInspectorChannel)

	var newLeader string
	for {
		select {
		case <-stateManagerChannel:
			e.hydraConfig.EtcdConf.Peers = []string{}
			e.EtcdManager.Restart(e.hydraConfig.EtcdConf)
		case newLeader = <-clusterInspectorChannel:
			e.hydraConfig.EtcdConf.Peers = []string{newLeader}
			e.EtcdManager.Restart(e.hydraConfig.EtcdConf)
		}
	}
}
