package peer

import (
	"fmt"
	"net"

	"github.com/ak1ra24/mrgobgpd/pkg/config"
	"github.com/ak1ra24/mrgobgpd/pkg/connection"
	"github.com/ak1ra24/mrgobgpd/pkg/event"
	"github.com/ak1ra24/mrgobgpd/pkg/packets"
	"github.com/ak1ra24/mrgobgpd/pkg/state"
	"github.com/rs/zerolog/log"
)

type Peer struct {
	State      state.State
	EventQueue event.EventQueue
	Config     config.Config
	TcpConn    net.Conn
}

func NewPeer(config config.Config) *Peer {
	state := state.Idle
	eventQueue := event.NewEventQueue()

	return &Peer{
		State:      state,
		EventQueue: *eventQueue,
		Config:     config,
	}
}

func (peer *Peer) Start() {
	log.Info().Msg("peer is started")
	peer.EventQueue.Enqueue(event.ManualStart)
}

func (peer *Peer) Next() {
	event := peer.EventQueue.Dequeue()
	log.Info().Msgf("event is occured, event=%v", event)
	peer.handleEvent(event)
	fmt.Println("Next peer State: ", peer.State)
	if peer.TcpConn != nil {
		switch peer.State {
		case state.OpenSent:
			msg, err := connection.GetMessage(peer.TcpConn, &packets.Message{
				Body: &packets.OpenMessage{},
			})
			if err != nil {
				panic(err)
			}
			peer.handleMessage(msg)
		case state.OpenConfirm:
			msg, err := connection.GetMessage(peer.TcpConn, &packets.Message{
				Body: &packets.MessageHeader{},
			})
			if err != nil {
				panic(err)
			}
			peer.handleMessage(msg)
		}
	}
}

func (peer *Peer) handleEvent(eventInfo event.Event) {
	switch peer.State {
	case state.Idle:
		switch eventInfo {
		case event.ManualStart:
			c, err := connection.Connect(&peer.Config)
			if err != nil {
				log.Fatal().Msg(err.Error())
			}
			peer.EventQueue.Enqueue(event.TcpConnectionConfirmed)
			peer.State = state.Connect
			peer.TcpConn = c.Conn
		default:
		}
	case state.Connect:
		switch eventInfo {
		case event.TcpConnectionConfirmed:
			localIP := net.ParseIP(peer.Config.LocalIP)
			openMsg := packets.NewOpenMessage(uint16(peer.Config.LocalAS), 100, localIP)
			connection.Send(peer.TcpConn, *openMsg)
			peer.State = state.OpenSent
		default:
		}
	case state.OpenSent:
		switch eventInfo {
		case event.BgpOpen:
			keepAliveMsg := packets.NewKeepAliveMessage()
			if err := connection.Send(peer.TcpConn, *keepAliveMsg); err != nil {
				panic(err)
			}
			peer.State = state.OpenConfirm
			fmt.Println("OpenSent End")
		default:
		}
	case state.OpenConfirm:
		switch eventInfo {
		case event.KeepAliveMsg:
			peer.State = state.Established
			peer.EventQueue.Enqueue(event.Established)
		}
	default:
	}
}

func (peer *Peer) handleMessage(msg *packets.Message) {
	switch msg.Header.Type {
	case packets.MsgOpen:
		peer.EventQueue.Enqueue(event.BgpOpen)
	case packets.MsgKeepAlive:
		peer.EventQueue.Enqueue(event.KeepAliveMsg)
	case packets.MsgUpdate:
		peer.EventQueue.Enqueue(event.UpdateMsg)
	}
}
