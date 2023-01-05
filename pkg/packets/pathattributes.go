package packets

// import (
// 	"encoding/binary"
// 	"errors"
// 	"net"
// 	"net/netip"
// 	"sort"
// )

// type Origin uint16

// const (
// 	Igp Origin = iota
// 	Egp
// 	Incomplete
// )

// type AsPath struct {
// 	AsSequence []int
// 	AsSet      sort.IntSlice
// }

// type Ipv4Network struct {
// 	Target net.IP
// 	Str    string
// }

// type PathAttribute struct {
// 	Origin   Origin
// 	AsPath   AsPath
// 	NextHop  Ipv4Network
// 	DontKnow []uint8
// }

// func (pa *PathAttribute) From() ([]byte, error) {
// 	var buf []byte
// 	if pa.Origin {
// 		attributeFlag := 0b01000000
// 		attributeTypeCode := 1
// 		attributeLength := 1
// 		attribute := 9999
// 		switch pa.Origin {
// 		case Igp:
// 			attribute = 0
// 		case Egp:
// 			attribute = 1
// 		case Incomplete:
// 			attribute = 2
// 		}

// 		buf[0] = byte(attributeFlag)
// 		buf[1] = byte(attributeTypeCode)
// 		buf[2] = byte(attributeLength)
// 		buf[3] = byte(attribute)
// 	}
// 	if pa.AsPath {
// 		attributeFlag := 0b01000000
// 		attributeTypeCode := 2
// 		attributeLength := pa.AsPath.BytesLen()
// 		attributeLengthBytes := []byte{}

// 		if attributeLength < 256 {
// 			attributeLengthBytes[0] = byte(attributeLength)
// 		} else {
// 			attributeFlag += 0b00010000
// 			binary.BigEndian.PutUint16(attributeLengthBytes, attributeLength)
// 		}

// 		attribute, _ := pa.AsPath.From()

// 		buf[0] = byte(attributeFlag)
// 		buf[1] = byte(attributeTypeCode)
// 		buf = append(buf, attributeLengthBytes...)
// 		buf = append(buf, attribute...)
// 	}
// 	if pa.NextHop {
// 		attributeFlag := 0b01000000
// 		attributeTypeCode := 3
// 		attributeLength := 4
// 		attribute, _ := pa.NextHop.Target.MarshalText()
// 		buf[0] = byte(attributeFlag)
// 		buf[1] = byte(attributeTypeCode)
// 		buf[2] = byte(attributeLength)
// 		buf = append(buf, attribute...)
// 	}
// 	// if pa.DontKnow {
// 	// 	buf = a
// 	// }

// 	return buf, nil
// }

// func (pa *PathAttribute) bytesLen() uint16 {
// 	var paLength uint16
// 	if pa.Origin != 9999 {
// 		paLength = 1
// 	}
// 	if (pa.AsPath != AsPath{}) {
// 		paLength = pa.AsPath.BytesLen()
// 	}
// 	if pa.NextHop != "" {
// 		paLength = 4
// 	}
// 	if pa.DontKnow != nil {
// 		paLength = uint16(len(pa.DontKnow))
// 	}

// 	length := paLength + 2

// 	if paLength > 255 {
// 		length += 2
// 	} else {
// 		length += 1
// 	}
// 	return length
// }

// func (ap *AsPath) BytesLen() uint16 {
// 	var asBytesLength uint16

// 	if ap.AsSequence != nil {
// 		asBytesLength = uint16(2 * len(ap.AsSequence))
// 	}
// 	if ap.AsSet != nil {
// 		asBytesLength = uint16(2 * len(ap.AsSet))
// 	}

// 	asBytesLength += 2

// 	return asBytesLength
// }

// func (ap *AsPath) From() ([]byte, error) {
// 	buf := []byte{}
// 	if ap.AsSet {
// 		pathSegmentType := 1
// 		numOfAses := len(ap.AsSet)
// 		buf[0] = byte(pathSegmentType)
// 		buf[1] = byte(numOfAses)
// 		for _, as := range ap.AsSet {
// 			buf = append(buf, byte(as))
// 		}
// 	}
// 	if ap.AsSequence {
// 		pathSegmentType := 2
// 		numOfAses := len(ap.AsSequence)
// 		buf[0] = byte(pathSegmentType)
// 		buf[1] = byte(numOfAses)
// 		for _, as := range ap.AsSequence {
// 			buf = append(buf, byte(as))
// 		}
// 	}
// 	return buf, nil
// }

// func (ip *Ipv4Network) BytesLen() (*uint16, error) {
// 	prefix, err := netip.ParsePrefix(ip.Str)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var usize uint16
// 	switch prefix.Bits() {
// 	case 0:
// 		usize = 1
// 		return &usize, nil
// 	case 1, 2, 3, 4, 5, 6, 7, 8:
// 		usize = 2
// 		return &usize, nil
// 	case 9, 10, 11, 12, 13, 14, 15, 16:
// 		usize = 3
// 		return &usize, nil
// 	case 17, 18, 19, 20, 21, 22, 23, 24:
// 		usize = 4
// 		return &usize, nil
// 	case 25, 26, 27, 28, 29, 30, 31, 32:
// 		usize = 5
// 		return &usize, nil
// 	default:
// 		return nil, errors.New("not prefix")
// 	}

// }
