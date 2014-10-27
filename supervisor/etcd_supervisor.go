package supervisor

import (
	"github.com/innotech/hydra/log"
	etcd_config "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"strings"
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
	addr := strings.Replace(config.Addr, "http://", "", 1)
	addr = strings.Replace(addr, "https://", "", 1)
	peerAddr := strings.Replace(config.Peer.Addr, "http://", "", 1)
	peerAddr = strings.Replace(peerAddr, "https://", "", 1)
	return &EtcdSupervisor{
		etcdConfig: config,
		// ClusterInspector: NewClusterInspector(config.Addr, config.Peer.Addr, configPeers),
		ClusterInspector: NewClusterInspector(addr, configPeers),
		EtcdManager:      NewEtcdManager(),
		StateManager: NewStateManager(
			time.Duration(config.Peer.HeartbeatTimeout*3)*time.Millisecond,
			NewEtcdClient([]string{config.Addr}).WithMachineAddr(config.Addr),
			addr,
			peerAddr,
			// TODO: should be calculates from first argument
			uint64(config.Peer.HeartbeatTimeout*6),
		),
	}
}

func (e *EtcdSupervisor) Run() {
	e.EtcdManager.Start(e.etcdConfig)

	stateManagerChannel := make(chan StateControllerState)
	go func() {
		e.StateManager.Run(stateManagerChannel)
	}()
	clusterInspectorChannel := make(chan string)
	go func() {
		e.ClusterInspector.Run(clusterInspectorChannel)
	}()

	var newLeader string
	for {
		select {
		case <-stateManagerChannel:
			log.Debug("Unable to write state")
			e.etcdConfig.Peers = []string{}
			e.EtcdManager.Restart(e.etcdConfig)
		case newLeader = <-clusterInspectorChannel:
			log.Debug("New leader found with peer address: " + newLeader)
			e.etcdConfig.Peers = []string{newLeader}
			e.EtcdManager.Restart(e.etcdConfig)
		}
	}
}
