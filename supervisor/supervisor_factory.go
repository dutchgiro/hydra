package supervisor

// import (
// 	hydra_config "github.com/innotech/hydra/config"

// 	"time"
// )

// const (
// 	StateMaster = iota
// 	StateMasterP
// 	Slave
// 	Unavailable
// )

// const (
// 	ClusterPath                  string = "/cluster/"
// 	StateKey                     string = "state"
// 	DefaultMaxAttemptsToSetState uint   = 3

// 	// DefaultAppsCacheDuration              time.Duration = time.Duration(20000) * time.Millisecond
// 	// DefaultDurationBetweenAllServersRetry time.Duration = time.Duration(0) * time.Millisecond
// 	// DefaultHydraServersCacheDuration      time.Duration = time.Duration(60000) * time.Millisecond
// 	// DefaultWriteAttempts uint = 3
// )

// type etcdSupervisorFactory struct {
// 	durationBetweenPublicationsState time.Duration
// 	peerAddr                         string
// 	stateTTL                         uint64

// 	// etcdAddr         string
// 	// maxWriteAttempts uint
// }

// var EtcdSupervisorFactory etcdSupervisorFactory = &etcdSupervisorFactory{
// 	maxWriteAttempts: DefaultWriteAttempts,
// }

// func (e *etcdSupervisorFactory) Config(config hydra_config.Config) {
// 	e.peerAddr = config.Peer.Addr
// 	millisecondsBetweenPublicationsState := config.Peer.HeartbeatTimeout * 2
// 	e.durationBetweenPublicationsState = time.Duration(millisecondsBetweenPublicationsState) * time.Millisecond
// 	e.stateTTL = millisecondsBetweenPublicationsState * (DefaultMaxAttemptsToSetState + 1)
// }

// func (e *etcdSupervisorFactory) Build() {
// 	supervisor := NewSupervisor(NewSupervisorRequester())
// 	supervisor.SetMaxWriteAttempts = e.maxWriteAttempts
// }

// // TODO: Validate

// func (e *etcdSupervisorFactory) GetEtcdAddr() uint {
// 	return e.maxWriteAttempts
// }

// func (e *etcdSupervisorFactory) GetMaxWriteAttempts() uint {
// 	return e.maxWriteAttempts
// }

// func (e *etcdSupervisorFactory) WithEtcdAddr(addr string) *supervisorFactory {
// 	e.etcdAddr = addr
// 	return e
// }

// func (e *etcdSupervisorFactory) WithMaxWriteAttempts(writeAttempts uint) *supervisorFactory {
// 	e.maxWriteAttempts = writeAttempts
// 	return e
// }
