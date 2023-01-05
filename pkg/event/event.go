package event

type Event int

const (
	ManualStart Event = iota
	TcpConnectionConfirmed
	BgpOpen
	KeepAliveMsg
	UpdateMsg
	Established
)
