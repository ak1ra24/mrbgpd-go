package config

import (
	"strconv"
	"strings"

	"github.com/ak1ra24/mrgobgpd/pkg/bgp"
)

type Config struct {
	LocalAS  bgp.AutonomousSystemNumber
	LocalIP  string
	RemoteAS bgp.AutonomousSystemNumber
	RemoteIP string
	Mode     Mode
	Networks []string
}

type Mode int

const (
	_ Mode = iota
	Passive
	Active
)

func ModeFromStr(s string) Mode {
	switch s {
	case "active":
		return Active
	case "passive":
		return Passive
	default:
		return 0
	}
}

func ParseConfig(configStr string) (*Config, error) {
	confStr := strings.Split(configStr, " ")
	localASNumber, err := strconv.ParseUint(confStr[0], 10, 16)
	if err != nil {
		return nil, err
	}
	remoteASNumber, err := strconv.ParseUint(confStr[2], 10, 16)
	if err != nil {
		return nil, err
	}
	mode := ModeFromStr(confStr[4])

	networks := []string{}
	for _, networkStr := range confStr[5:] {
		networks = append(networks, networkStr)
	}

	return &Config{
		LocalAS:  bgp.AutonomousSystemNumber(localASNumber),
		LocalIP:  confStr[1],
		RemoteAS: bgp.AutonomousSystemNumber(remoteASNumber),
		RemoteIP: confStr[3],
		Mode:     mode,
		Networks: networks,
	}, nil
}
