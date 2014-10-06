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
	DefaultMaxAttemptsToSetState uint   = 2
	StateAvailable               string = "available"
	StateKey                     string = "state"
)

type EtcdSupervisor struct {
	// EtcdClient *EtcdClient
	EtcdClient EtcdRequester

	durationBetweenPublicationsState time.Duration
	modifiedIndexKeyState            uint64
	numOfSetStateRetries             uint
	peerAddr                         string
	state                            uint
	stateTTL                         uint64

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

// func NewEtcdSupervisor(config *hydra_config.Config) *EtcdSupervisor {
// 	millisecondsBetweenPublicationsState := uint(config.Peer.HeartbeatTimeout * 2)

// 	return &EtcdSupervisor{
// 		EtcdClient:                       NewEtcdClient([]string{config.EtcdAddr}),
// 		durationBetweenPublicationsState: time.Duration(millisecondsBetweenPublicationsState) * time.Millisecond,
// 		modifiedIndexKeyState:            0,
// 		numOfSetStateRetries:             DefaultMaxAttemptsToSetState,
// 		peerAddr:                         config.Peer.Addr,
// 		stateTTL:                         uint64(millisecondsBetweenPublicationsState * (DefaultMaxAttemptsToSetState + 1)),
// 	}
// }

func NewEtcdSupervisor(config *hydra_config.Config, etcdClient EtcdRequester) *EtcdSupervisor {
	millisecondsBetweenPublicationsState := uint(config.Peer.HeartbeatTimeout * 2)

	return &EtcdSupervisor{
		EtcdClient:                       etcdClient,
		durationBetweenPublicationsState: time.Duration(millisecondsBetweenPublicationsState) * time.Millisecond,
		modifiedIndexKeyState:            0,
		numOfSetStateRetries:             DefaultMaxAttemptsToSetState,
		peerAddr:                         config.Peer.Addr,
		stateTTL:                         uint64(millisecondsBetweenPublicationsState * (DefaultMaxAttemptsToSetState + 1)),
	}
}

func (e *EtcdSupervisor) Run() {
	_ = e.setNodeState()
}

func (e *EtcdSupervisor) setNodeState() error {
	fmt.Println("-------> setNodeState")
	var res *Response
	var err error = nil
	stateKey := ClusterRootPath + "/" + e.peerAddr + "/" + StateKey
	// fmt.Println(stateKey)
	// fmt.Println(StateAvailable)
	// fmt.Printf("%d\n", e.stateTTL)

	if e.modifiedIndexKeyState > 0 {
		for i := 0; i < int(e.numOfSetStateRetries); i++ {
			res, err = e.EtcdClient.CompareAndSwap(stateKey, StateAvailable, e.stateTTL, "", e.modifiedIndexKeyState)
			if err == nil {
				e.modifiedIndexKeyState = res.Node.ModifiedIndex
				break
			}
		}
	} else {
		fmt.Printf("-------> First: %d\n", e.numOfSetStateRetries)

		for i := 0; i < int(e.numOfSetStateRetries); i++ {
			fmt.Println("-------> Set")
			fmt.Println(stateKey)
			fmt.Println(StateAvailable)
			fmt.Printf("%d\n", e.stateTTL)
			res, err = e.EtcdClient.Set(stateKey, StateAvailable, e.stateTTL)
			fmt.Printf("%#v", res)
			fmt.Printf("%#v\n", err)
			if res != nil {
				// e.modifiedIndexKeyState = res.Node.ModifiedIndex
				// break
			}
		}
	}

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
