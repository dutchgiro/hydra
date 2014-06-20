package load_balancer

const (
	//  Server commands, as strings
	SIGNAL_READY      = "\001"
	SIGNAL_REQUEST    = "\002"
	SIGNAL_REPLY      = "\003"
	SIGNAL_HEARTBEAT  = "\004"
	SIGNAL_DISCONNECT = "\005"
)

var Commands = []string{"", "READY", "REQUEST", "REPLY", "HEARTBEAT", "DISCONNECT"}
