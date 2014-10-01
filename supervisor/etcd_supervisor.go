package supervisor

import (
	"time"
)

type Supervisor interface {
	Run()
}

type EtcdSupervisor struct {
	etcdRequester                    *EtcdRequester
	maxAttemptsToSetState            uint
	durationBetweenPublicationsState time.Duration
	peerAddr                         string
	state                            uint
}

func newEtcdSupervisor(requester *EtcdRequester) {
	return &Supervisor{
		etcdRequester: requester,
	}
}

func (s *Supervisor) Run() {

}

func (s *Supervisor) setMaxAttemptsToSetState(attempts uint) {
	s.maxAttemptsToSetState = attempts
}

func (s *Supervisor) GetDurationBetweenPublicationsState() time.Duration {
	return durationBetweenPublicationsState
}

func (s *Supervisor) GetMaxAttemptsToSetState() uint {
	return s.maxAttemptsToSetState
}

func (s *Supervisor) GetPeerAddr() string {
	return s.peerAddr
}

func (s *Supervisor) GetState() uint {
	return s.state
}
