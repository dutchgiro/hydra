package supervisor

import (
	hydra_config "github.com/innotech/hydra/config"

	// "fmt"
	"time"
)

type Supervisor interface {
	Run()
}

// const (
// 	ClusterRootPath              string = "cluster"
// 	DefaultMaxAttemptsToSetState uint   = 3
// 	// StateAvailable               string = "available"
// 	StateKey string = "state"
// )

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

////////////////////////////////////////////////////////////////////////////

// type EtcdSupervisor struct {
// 	EtcdClient  EtcdRequester
// 	EtcdManager EtcdController

// 	durationBetweenPublicationsState time.Duration
// 	modifiedIndexKeyState            uint64
// 	numOfSetStateRetries             uint
// 	peerAddr                         string
// 	state                            uint
// 	stateTTL                         uint64
// 	executionChannel                 chan int
// }

// func NewEtcdSupervisor(config *hydra_config.Config) *EtcdSupervisor {
// 	millisecondsBetweenPublicationsState := uint(config.Peer.HeartbeatTimeout * 2)

// 	return &EtcdSupervisor{
// 		EtcdClient:                       NewEtcdClient([]string{config.EtcdAddr}),
// 		durationBetweenPublicationsState: time.Duration(millisecondsBetweenPublicationsState) * time.Millisecond,
// 		executionChannel:                 make(chan int),
// 		modifiedIndexKeyState:            0,
// 		numOfSetStateRetries:             DefaultMaxAttemptsToSetState,
// 		peerAddr:                         config.Peer.Addr,
// 		stateTTL:                         uint64(millisecondsBetweenPublicationsState * (DefaultMaxAttemptsToSetState + 1)),
// 	}
// }

// func (e *EtcdSupervisor) Run() {
// 	fmt.Println("RUN")

// OuterLoop:
// 	for {
// 		fmt.Println("------->>> FOR")
// 		select {
// 		case <-e.executionChannel:
// 			fmt.Println("------->>> BREAK LOOP")
// 			break OuterLoop
// 		default:
// 			fmt.Println("------->>> Call to setNodeState")
// 			err := e.setNodeState()
// 			if err != nil {
// 				fmt.Println("------->>> Call to Restart")
// 				e.EtcdManager.Restart()
// 				fmt.Println("------->>> Break")
// 				// break OuterLoop
// 			}
// 			time.Sleep(e.durationBetweenPublicationsState)
// 		}
// 	}

// 	fmt.Println("EXIT RUN")
// }

// func (e *EtcdSupervisor) Stop() {
// 	e.executionChannel <- 0
// }

// func (e *EtcdSupervisor) setNodeState() error {
// 	var res *Response
// 	var err error = nil
// 	stateKey := ClusterRootPath + "/" + e.peerAddr + "/" + StateKey

// 	var stateKeyExistence KeyExistence = False
// 	if e.modifiedIndexKeyState > 0 {
// 		stateKeyExistence = Unknow
// 	}

// 	for i := 0; i < int(e.numOfSetStateRetries); i++ {
// 		fmt.Println("*************")
// 		res, err = e.EtcdClient.CompareAndSwap(stateKey, StateAvailable, e.stateTTL,
// 			"", e.modifiedIndexKeyState, stateKeyExistence)
// 		// res, err = e.EtcdClient.CompareAndSwap(stateKey, StateAvailable, e.stateTTL,
// 		// 	"", e.modifiedIndexKeyState)
// 		fmt.Println("*************")
// 		if err == nil && res != nil {
// 			e.modifiedIndexKeyState = res.Node.ModifiedIndex
// 			break
// 		}
// 		time.Sleep(e.durationBetweenPublicationsState)
// 	}

// 	return err
// }

// // func (e *Supervisor) setMaxAttemptsToSetState(attempts uint) {
// // 	e.maxAttemptsToSetState = attempts
// // }

// // func (e *Supervisor) setPeerAddr(peerAddr uint) {
// // 	e.peerAddr = peerAddr
// // }

// func (e *EtcdSupervisor) GetNumOfSetStateRetries() uint {
// 	return e.numOfSetStateRetries
// }

// // func (e *Supervisor) GetMaxAttemptsToSetState() uint {
// // 	return e.maxAttemptsToSetState
// // }

// func (e *EtcdSupervisor) GetPeerAddr() string {
// 	return e.peerAddr
// }

// func (e *EtcdSupervisor) GetState() uint {
// 	return e.state
// }

// func (e *EtcdSupervisor) GetStateTTL() uint64 {
// 	return e.stateTTL
// }
