package supervisor

type EtcdController interface {
	Restart()
	Start()
	Stop()
}

type EtcdManager struct {
}

func (e *EtcdManager) Restart() {

}

func (e *EtcdManager) Start() {

}

func (e *EtcdManager) Stop() {

}
