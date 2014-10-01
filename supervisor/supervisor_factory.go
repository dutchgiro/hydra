package supervisor

const (
	StateMaster = iota
	StateMasterP
	Slave
	Unavailable
)

const (
	ClusterPath string = "/cluster/"
	StateKey    string = "state"

	// DefaultAppsCacheDuration              time.Duration = time.Duration(20000) * time.Millisecond
	// DefaultDurationBetweenAllServersRetry time.Duration = time.Duration(0) * time.Millisecond
	// DefaultHydraServersCacheDuration      time.Duration = time.Duration(60000) * time.Millisecond
	DefaultWriteAttempts uint = 3
)

type etcdSupervisorFactory struct {
	etcdAddr         string
	maxWriteAttempts uint
}

var EtcdSupervisorFactory etcdSupervisorFactory = &etcdSupervisorFactory{
	maxWriteAttempts: DefaultWriteAttempts,
}

func (e *etcdSupervisorFactory) Config(etcdAddr string) {
	e.etcdAddr = etcdAddr
}

func (e *etcdSupervisorFactory) Build() {
	supervisor := NewSupervisor(NewSupervisorRequester())
	supervisor.SetMaxWriteAttempts = e.maxWriteAttempts
}

// TODO: Validate

func (e *etcdSupervisorFactory) GetEtcdAddr() uint {
	return e.maxWriteAttempts
}

func (e *etcdSupervisorFactory) GetMaxWriteAttempts() uint {
	return e.maxWriteAttempts
}

func (e *etcdSupervisorFactory) WithEtcdAddr(addr string) *supervisorFactory {
	e.etcdAddr = addr
	return e
}

func (e *etcdSupervisorFactory) WithMaxWriteAttempts(writeAttempts uint) *supervisorFactory {
	e.maxWriteAttempts = writeAttempts
	return e
}
