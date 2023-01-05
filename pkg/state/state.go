package state

type State int

const (
	Idle State = iota
	Connect
	OpenSent
	OpenConfirm
	Established
)
