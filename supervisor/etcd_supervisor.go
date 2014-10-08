package supervisor

import (
	hydra_config "github.com/innotech/hydra/config"

	"fmt"
	"time"
)

type Supervisor interface {
	Run()
}

const (
	ClusterRootPath              string = "cluster"
	DefaultMaxAttemptsToSetState uint   = 3
	StateAvailable               string = "available"
	StateKey                     string = "state"
)

type EtcdSupervisor struct {
	EtcdClient  EtcdRequester
	EtcdManager EtcdController

	durationBetweenPublicationsState time.Duration
	modifiedIndexKeyState            uint64
	numOfSetStateRetries             uint
	peerAddr                         string
	state                            uint
	stateTTL                         uint64
	executionChannel                 chan int

	// etcdRequester                    *EtcdRequester
	// maxAttemptsToSetState            uint
	// durationBetweenPublicationsState time.Duration
	// peerAddr                         string
	// state                            uint
}

// func newEtcdSupervisor(peerAddr string, durationBetweenPublicationsState time.Duration) {
// 	return &Supervisor{
// 		durationBetweenPublicationsState: durationBetweenPublicationsState,
// 		peerAddr: peerAddr,
// 	}
// }

func NewEtcdSupervisor(config *hydra_config.Config) *EtcdSupervisor {
	millisecondsBetweenPublicationsState := uint(config.Peer.HeartbeatTimeout * 2)

	return &EtcdSupervisor{
		EtcdClient:                       NewEtcdClient([]string{config.EtcdAddr}),
		durationBetweenPublicationsState: time.Duration(millisecondsBetweenPublicationsState) * time.Millisecond,
		executionChannel:                 make(chan int),
		modifiedIndexKeyState:            0,
		numOfSetStateRetries:             DefaultMaxAttemptsToSetState,
		peerAddr:                         config.Peer.Addr,
		stateTTL:                         uint64(millisecondsBetweenPublicationsState * (DefaultMaxAttemptsToSetState + 1)),
	}
}

// func NewEtcdSupervisor(config *hydra_config.Config, etcdClient EtcdRequester) *EtcdSupervisor {
// 	millisecondsBetweenPublicationsState := uint(config.Peer.HeartbeatTimeout * 2)

// 	return &EtcdSupervisor{
// 		EtcdClient:                       etcdClient,
// 		durationBetweenPublicationsState: time.Duration(millisecondsBetweenPublicationsState) * time.Millisecond,
// 		modifiedIndexKeyState:            0,
// 		numOfSetStateRetries:             DefaultMaxAttemptsToSetState,
// 		peerAddr:                         config.Peer.Addr,
// 		stateTTL:                         uint64(millisecondsBetweenPublicationsState * (DefaultMaxAttemptsToSetState + 1)),
// 	}
// }

func (e *EtcdSupervisor) Run() {
	fmt.Println("RUN")
OuterLoop:
	for {
		fmt.Println("------->>> FOR")
		select {
		case <-e.executionChannel:
			fmt.Println("------->>> BREAK LOOP")
			break OuterLoop
		default:
			fmt.Println("------->>> Call to setNodeState")
			err := e.setNodeState()
			if err != nil {
				fmt.Println("------->>> Call to Restart")
				e.EtcdManager.Restart()
				fmt.Println("------->>> Break")
				// break OuterLoop
			}
			time.Sleep(e.durationBetweenPublicationsState)
		}
	}

	fmt.Println("EXIT RUN")
}

func (e *EtcdSupervisor) Stop() {
	e.executionChannel <- 0
}

func (e *EtcdSupervisor) setNodeState() error {
	var res *Response
	var err error = nil
	stateKey := ClusterRootPath + "/" + e.peerAddr + "/" + StateKey

	var stateKeyExistence KeyExistence = False
	if e.modifiedIndexKeyState > 0 {
		stateKeyExistence = Unknow
	}

	for i := 0; i < int(e.numOfSetStateRetries); i++ {
		res, err = e.EtcdClient.CompareAndSwap(stateKey, StateAvailable, e.stateTTL,
			"", e.modifiedIndexKeyState, stateKeyExistence)
		if err == nil && res != nil {
			e.modifiedIndexKeyState = res.Node.ModifiedIndex
			break
		}
		time.Sleep(e.durationBetweenPublicationsState)
	}

	// if e.modifiedIndexKeyState > 0 {
	// 	fmt.Println("------->>> CompareAndSwap")
	// 	fmt.Printf("%d\n", e.numOfSetStateRetries)
	// 	for i := 0; i < int(e.numOfSetStateRetries); i++ {
	// 		fmt.Println(stateKey)
	// 		fmt.Println(StateAvailable)
	// 		fmt.Printf("%d\n", e.stateTTL)
	// 		fmt.Printf("%d\n", e.modifiedIndexKeyState)
	// 		res, err = e.EtcdClient.CompareAndSwap(stateKey, StateAvailable, e.stateTTL, "", e.modifiedIndexKeyState)
	// 		if err == nil && res != nil {
	// 			e.modifiedIndexKeyState = res.Node.ModifiedIndex
	// 			break
	// 		}
	// 		time.Sleep(e.durationBetweenPublicationsState)
	// 	}
	// } else {
	// 	fmt.Println("------->>> Set")
	// 	fmt.Printf("%d\n", e.numOfSetStateRetries)
	// 	for i := 0; i < int(e.numOfSetStateRetries); i++ {
	// 		fmt.Println(stateKey)
	// 		fmt.Println(StateAvailable)
	// 		fmt.Printf("%d\n", e.stateTTL)
	// 		res, err = e.EtcdClient.Set(stateKey, StateAvailable, e.stateTTL)
	// 		// fmt.Printf("%#v\n", res)
	// 		// fmt.Printf("%#v\n", err)
	// 		if err == nil && res != nil {
	// 			// fmt.Printf("%d\n", res.Node.ModifiedIndex)
	// 			e.modifiedIndexKeyState = res.Node.ModifiedIndex
	// 			break
	// 		}
	// 		time.Sleep(e.durationBetweenPublicationsState)
	// 	}
	// }

	return err
}

// func (e *Supervisor) setMaxAttemptsToSetState(attempts uint) {
// 	e.maxAttemptsToSetState = attempts
// }

// func (e *Supervisor) setPeerAddr(peerAddr uint) {
// 	e.peerAddr = peerAddr
// }

func (e *EtcdSupervisor) GetNumOfSetStateRetries() uint {
	return e.numOfSetStateRetries
}

// func (e *Supervisor) GetMaxAttemptsToSetState() uint {
// 	return e.maxAttemptsToSetState
// }

func (e *EtcdSupervisor) GetPeerAddr() string {
	return e.peerAddr
}

func (e *EtcdSupervisor) GetState() uint {
	return e.state
}

func (e *EtcdSupervisor) GetStateTTL() uint64 {
	return e.stateTTL
}
