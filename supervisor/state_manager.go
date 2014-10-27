package supervisor

import (
	"github.com/innotech/hydra/log"

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

const (
	PersistentTTL uint64 = 0
)

// TODO: stateManager not exporting
type StateManager struct {
	DurationBetweenPublicationsState time.Duration
	etcdClient                       EtcdRequester
	modifiedIndexKeyState            uint64
	NumOfSetStateRetries             uint
	PeerAddr                         string
	PeerAddrKey                      string
	PeerStateKey                     string
	PeerStateKeyTTL                  uint64
	state                            StateControllerState
}

func NewStateManager(durationBetweenPublicationsState time.Duration, etcdClient EtcdRequester,
	addr string, peerAddr string, stateKeyTTL uint64) *StateManager {
	return &StateManager{
		DurationBetweenPublicationsState: durationBetweenPublicationsState,
		etcdClient:                       etcdClient,
		modifiedIndexKeyState:            0,
		NumOfSetStateRetries:             DefaultMaxAttemptsToSetState,
		PeerAddr:                         peerAddr,
		PeerAddrKey:                      ClusterKey + "/" + addr + "/" + PeerAddrKey,
		PeerStateKey:                     ClusterKey + "/" + addr + "/" + StateKey,
		PeerStateKeyTTL:                  stateKeyTTL,
		state:                            Stopped,
	}
}

func (s *StateManager) GetState() StateControllerState {
	return s.state
}

func (s *StateManager) Run(ch chan StateControllerState) {
	s.state = Running
	_, err := s.etcdClient.Set(s.PeerAddrKey, s.PeerAddr, PersistentTTL)
	if err != nil {
		log.Warn("Unable to set peer address in cluster key")
	}
	log.Debug("Peer address was set with value: " + s.PeerAddr)
OuterLoop:
	for {
		err := s.setPeerState()
		if err != nil {
			s.state = Stopped
			ch <- s.state
			break OuterLoop
		}
		log.Debug("Peer state " + PeerStateEnabled + " posted successfully")
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

func (s *StateManager) setPeerState() error {
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
