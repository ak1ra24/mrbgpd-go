package connection

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/ak1ra24/mrgobgpd/pkg/config"
	"github.com/ak1ra24/mrgobgpd/pkg/packets"
)

const BGPPort = 179

type Connection struct {
	Conn   net.Conn
	Buffer []byte
}

func Connect(cfg *config.Config) (*Connection, error) {
	var err error
	var conn net.Conn
	switch cfg.Mode {
	case config.Active:
		addr := fmt.Sprintf("%s:%d", cfg.RemoteIP, BGPPort)
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
	case config.Passive:
		addr := fmt.Sprintf("%s:%d", cfg.LocalIP, BGPPort)
		listner, err := net.Listen("tcp", addr)
		if err != nil {
			return nil, err
		}
		conn, err = listner.Accept()
		if err != nil {
			return nil, err
		}
	default:
	}

	return &Connection{
		Conn:   conn,
		Buffer: []byte{},
	}, nil
}

func Send(conn net.Conn, msg packets.Message) error {
	var data []byte
	hdr, err := msg.Header.From()
	if err != nil {
		return err
	}
	data = append(data, hdr...)

	if msg.Body != nil {
		body, err := msg.Body.From()
		if err != nil {
			return err
		}
		data = append(data, body...)
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func GetMessage(conn net.Conn, msg *packets.Message) (*packets.Message, error) {
	buf, err := readDataFromTcpConnection(conn)
	if err != nil {
		return nil, err
	}
	b, err := splitBufferAtMessageSeparator(buf)
	if err != nil {
		return nil, err
	}

	if err := msg.Header.TryFrom(b); err != nil {
		return nil, err
	}

	if err := msg.Body.TryFrom(b); err != nil {
		return nil, err
	}

	return msg, nil
}

func readDataFromTcpConnection(conn net.Conn) ([]byte, error) {
	var buf = make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func splitBufferAtMessageSeparator(buf []byte) ([]byte, error) {
	index, err := getIndexOfMesasgeSeparator(buf)
	if err != nil {
		return nil, err
	}
	if len(buf) < int(index) {
		return nil, errors.New("No Receive")
	}
	return buf[index:], nil
}

// bufferのうちどこまでが1つのbgp messageを表すbytesであるかを示す
func getIndexOfMesasgeSeparator(buf []byte) (uint16, error) {
	if len(buf) < packets.HEADER_LENGTH {
		return 0, errors.New("data is too short")
	}
	usize := binary.BigEndian.Uint16(buf[16:18])
	return usize, nil
}
