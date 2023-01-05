package peer

import (
	"log"
	"testing"
	"time"

	"github.com/ak1ra24/mrgobgpd/pkg/config"
	"github.com/ak1ra24/mrgobgpd/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestPeerCanTransitionToConnectState(t *testing.T) {
	configStr := "64512 127.0.0.1 64513 127.0.0.2 active"
	cfg, err := config.ParseConfig(configStr)
	if err != nil {
		t.Fatal(err)
	}
	peer := NewPeer(*cfg)
	peer.Start()

	go func() {
		remoteConfigStr := "64513 127.0.0.2 64512 127.0.0.1 passive"
		remoteCfg, err := config.ParseConfig(remoteConfigStr)
		if err != nil {
			log.Println(err)
		}
		remotePeer := NewPeer(*remoteCfg)
		remotePeer.Start()
		remotePeer.Next()
	}()

	time.Sleep(10 * time.Second)

	peer.Next()

	assert.Equal(t, peer.State, state.Connect)
}

func TestPeerCanTransitionToOpenSentState(t *testing.T) {
	configStr := "64512 127.0.0.1 64513 127.0.0.2 active"
	cfg, err := config.ParseConfig(configStr)
	if err != nil {
		t.Fatal(err)
	}
	peer := NewPeer(*cfg)
	peer.Start()

	go func() {
		remoteConfigStr := "64513 127.0.0.2 64512 127.0.0.1 passive"
		remoteCfg, err := config.ParseConfig(remoteConfigStr)
		if err != nil {
			log.Println(err)
		}
		remotePeer := NewPeer(*remoteCfg)
		remotePeer.Start()
		remotePeer.Next()
		remotePeer.Next()
	}()

	time.Sleep(10 * time.Second)

	peer.Next()
	peer.Next()

	assert.Equal(t, peer.State, state.OpenSent)
}

func TestPeerCanTransitionToOpenConfirm(t *testing.T) {
	configStr := "64512 127.0.0.1 64513 127.0.0.2 active"
	cfg, err := config.ParseConfig(configStr)
	if err != nil {
		t.Fatal(err)
	}
	peer := NewPeer(*cfg)
	peer.Start()

	go func() {
		remoteConfigStr := "64513 127.0.0.2 64512 127.0.0.1 passive"
		remoteCfg, err := config.ParseConfig(remoteConfigStr)
		if err != nil {
			log.Println(err)
		}
		remotePeer := NewPeer(*remoteCfg)
		remotePeer.Start()
		maxStep := 50
		for i := 1; i <= maxStep; i++ {
			remotePeer.Next()
			if remotePeer.State == state.OpenConfirm {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(10 * time.Second)

	maxStep := 50
	for i := 1; i <= maxStep; i++ {
		peer.Next()
		if peer.State == state.OpenConfirm {
			break
		}
		time.Sleep(1 * time.Second)
	}

	assert.Equal(t, peer.State, state.OpenConfirm)
}

func TestPeerCanTransitionToEstablishedState(t *testing.T) {
	configStr := "64512 127.0.0.1 64513 127.0.0.2 active"
	cfg, err := config.ParseConfig(configStr)
	if err != nil {
		t.Fatal(err)
	}
	peer := NewPeer(*cfg)
	peer.Start()

	go func() {
		remoteConfigStr := "64513 127.0.0.2 64512 127.0.0.1 passive"
		remoteCfg, err := config.ParseConfig(remoteConfigStr)
		if err != nil {
			log.Println(err)
		}
		remotePeer := NewPeer(*remoteCfg)
		remotePeer.Start()
		maxStep := 50
		for i := 1; i <= maxStep; i++ {
			remotePeer.Next()
			if remotePeer.State == state.Established {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(10 * time.Second)

	maxStep := 50
	for i := 1; i <= maxStep; i++ {
		peer.Next()
		if peer.State == state.Established {
			break
		}
		time.Sleep(1 * time.Second)
	}

	assert.Equal(t, peer.State, state.Established)
}
