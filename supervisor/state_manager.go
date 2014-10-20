package supervisor

import (
	"time"
)

type StateController interface {
	Run(ch chan StateControllerState)
}

// const (
// 	ClusterRootPath              string = "cluster"
// 	DefaultMaxAttemptsToSetState uint   = 3
// 	StateAvailable               string = "available"
// 	StateKey                     string = "state"
// )

// const (
// 	PeerStateAvailable
// )

type StateControllerState uint8

const (
	Stopped StateControllerState = iota
	Running
)

// TODO: stateManager not exporting
type StateManager struct {
	DurationBetweenPublicationsState time.Duration
	etcdClient                       EtcdRequester
	modifiedIndexKeyState            uint64
	NumOfSetStateRetries             uint
	PeerStateKey                     string
	PeerStateKeyTTL                  uint64
	state                            StateControllerState
}

func NewStateManager(durationBetweenPublicationsState time.Duration, etcdClient EtcdRequester,
	peerAddr string, stateKeyTTL uint64) *StateManager {
	return &StateManager{
		DurationBetweenPublicationsState: durationBetweenPublicationsState,
		etcdClient:                       etcdClient,
		modifiedIndexKeyState:            0,
		NumOfSetStateRetries:             DefaultMaxAttemptsToSetState,
		PeerStateKey:                     ClusterKey + "/" + peerAddr + "/" + StateKey,
		PeerStateKeyTTL:                  stateKeyTTL,
		state:                            Stopped,
	}
}

func (s *StateManager) GetState() StateControllerState {
	return s.state
}

func (s *StateManager) Run(ch chan StateControllerState) {
	s.state = Running
OuterLoop:
	for {
		err := s.setNodeState()
		if err != nil {
			s.state = Stopped
			ch <- s.state
			break OuterLoop
		}
		time.Sleep(s.DurationBetweenPublicationsState)
	}
}

func (s *StateManager) getStateKeyExistence() (stateKeyExistence KeyExistence) {
	stateKeyExistence = False
	if s.modifiedIndexKeyState > 0 {
		stateKeyExistence = Unknow
	}
	return
}

func (s *StateManager) setNodeState() error {
	var res *Response
	var err error = nil

	stateKeyExistence := s.getStateKeyExistence()
	for i := 0; i < int(s.NumOfSetStateRetries); i++ {
		res, err = s.etcdClient.CompareAndSwap(s.PeerStateKey, PeerStateEnabled, s.PeerStateKeyTTL,
			"", s.modifiedIndexKeyState, stateKeyExistence)
		if err == nil && res != nil {
			s.modifiedIndexKeyState = res.Node.ModifiedIndex
			break
		}
		time.Sleep(s.DurationBetweenPublicationsState)
	}

	return err
}
