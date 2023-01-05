package packets

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	_ = iota
	Open
	KeepAlive
	Update
)

const (
	HEADER_LENGTH      = 19
	MAX_MESSAGE_LENGTH = 4096
)

type Message struct {
	Header MessageHeader
	Body   MessageBody
}

type MessageHeader struct {
	Marker []byte
	Length uint16
	Type   MessageType
}

type MessageBody interface {
	TryFrom([]byte) error
	From() ([]byte, error)
}

type MessageType uint8

const (
	_ MessageType = iota
	MsgOpen
	MsgKeepAlive
	MsgUpdate
)

func (mh *MessageHeader) From() ([]byte, error) {
	buf := make([]byte, HEADER_LENGTH)
	// BGPコネクションに認証を使用する場合に使う。もし、認証を使用しない場合はすべて"1"が入る。
	for i := range buf[:16] {
		buf[i] = 0xff
	}
	// ヘッダを含むメッセージ全体の長さをオクテット単位で定義されます。フィールドの値は最小で19byte、最大で4096byte
	binary.BigEndian.PutUint16(buf[16:18], mh.Length)
	// メッセージのタイプコードが定義されます。
	//1：OPEN
	//2：UPDATE
	//3：NOTIFICATION
	//4：KEEPALIVE
	buf[18] = byte(mh.Type)
	return buf, nil
}

func (mh *MessageHeader) TryFrom(data []byte) error {
	if uint16(len(data)) < HEADER_LENGTH {
		return fmt.Errorf("not all BGP message header")
	}
	// mh.Length = binary.BigEndian.Uint16(data[16:18])
	// if int(mh.Length) < HEADER_LENGTH {
	// 	return fmt.Errorf("unknown message type")
	// }
	mh.Type = MessageType(data[18])
	return nil
}

type OpenMessage struct {
	Version                 uint8
	AS                      uint16
	HoldTime                uint16
	BgpIdentifier           net.IP
	OptionalParameterLength uint8
	OptionalParameters      []uint8
}

func NewOpenMessage(as uint16, holdtime uint16, bgpIdentifier net.IP) *Message {
	return &Message{
		Header: MessageHeader{
			Type: MsgOpen,
		},
		Body: &OpenMessage{
			Version:                 4,
			AS:                      as,
			HoldTime:                holdtime,
			BgpIdentifier:           bgpIdentifier,
			OptionalParameterLength: 0,
			OptionalParameters:      []uint8{},
		},
	}
}

func (om *OpenMessage) TryFrom(data []byte) error {
	if len(data) < 10 {
		// 今はとりあえずのエラーメッセージ
		return fmt.Errorf("not all BGP Open message bytes available 1")
	}
	om.Version = data[0]
	om.AS = binary.BigEndian.Uint16(data[1:3])
	om.HoldTime = binary.BigEndian.Uint16(data[3:5])
	om.BgpIdentifier = net.IP(data[5:9]).To4()
	om.OptionalParameterLength = data[9]
	data = data[10:]
	// if len(data) < int(om.OptionalParameterLength) {
	// 	return fmt.Errorf("not all BGP Open message bytes available 2")
	// }
	// optional parameterの処理が入るが省略
	return nil
}

func (om *OpenMessage) From() ([]byte, error) {
	// versionからOptparameterLengthまでをバッファに入れる
	buf := make([]byte, 10)
	buf[0] = om.Version
	binary.BigEndian.PutUint16(buf[1:3], om.AS)
	binary.BigEndian.PutUint16(buf[3:5], om.HoldTime)
	copy(buf[5:9], om.BgpIdentifier.To4()) // 4byte分をコピー
	// optparamは現在使用していないので、そのまま何もせずにlengthを返す
	om.OptionalParameterLength = uint8(len(om.OptionalParameters))
	buf[9] = om.OptionalParameterLength
	return append(buf, om.OptionalParameters...), nil
}

func NewKeepAliveMessage() *Message {
	return &Message{
		Header: MessageHeader{
			Type: MsgKeepAlive,
		},
		Body: nil,
	}
}

// type UpdateMessage struct {
// 	WithDrawnRoutes                     []Ipv4Network
// 	WithDrawnRoutesLength               uint16
// 	PathAttributes                      []PathAttribute
// 	PathAttributesLength                uint16
// 	NetworkLayerReachabilityInformation []Ipv4Network
// }

// func NewUpdateMessage(pathAttributes []PathAttribute, nlri []Ipv4Network, withdrawnRoutes []Ipv4Network) *Message {
// 	PathAttributesLength := len(pathAttributes)
// 	nlriLength := len(nlri)
// 	withdrawnRoutesLength := len(withdrawnRoutes)
// 	return &Message{
// 		Header: MessageHeader{
// 			Length: uint16(HEADER_LENGTH + PathAttributesLength + nlriLength + withdrawnRoutesLength + 4),
// 			Type:   MsgUpdate,
// 		},
// 		Body: &UpdateMessage{
// 			WithDrawnRoutes:                     withdrawnRoutes,
// 			WithDrawnRoutesLength:               uint16(withdrawnRoutesLength),
// 			PathAttributes:                      pathAttributes,
// 			PathAttributesLength:                uint16(PathAttributesLength),
// 			NetworkLayerReachabilityInformation: nlri,
// 		},
// 	}
// }

// func (um *UpdateMessage) From() ([]byte, error) {
// 	var buf []byte
// 	binary.BigEndian.PutUint16(buf[0:2], um.WithDrawnRoutesLength)
// 	for _, withdrawnRoute := range um.WithDrawnRoutes {
// 		binary.BigEndian.AppendUint16(buf, withdrawnRoute.into())
// 	}

// 	binary.BigEndian.PutUint16(buf, um.PathAttributesLength)

// 	for _, pathAttribute := range um.PathAttributes {
// 		binary.BigEndian.AppendUint16(buf, pathAttribute.into())
// 	}

// 	for _, nlri := range um.NetworkLayerReachabilityInformation {
// 		binary.BigEndian.AppendUint16(buf, nlri)
// 	}
// }

// func (um *UpdateMessage) TryFrom(data []byte) error {

// }
