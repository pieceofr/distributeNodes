package p2p

import (
	"github.com/bitmark-inc/logger"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
)

var log *logger.L

// P2P structure holds information on currently running streams/Listeners
type P2P struct {
	ListenersLocal *Listeners
	ListenersP2P   *Listeners
	Streams        *StreamRegistry

	identity  peer.ID
	peerHost  *p2phost.Host
	peerstore pstore.Peerstore
}

// New creates new P2P struct
func New(identity peer.ID, peerHost *p2phost.Host, peerstore pstore.Peerstore, logger *logger.L) *P2P {
	log = logger
	log.Info("P2P: create a new P2P")
	return &P2P{
		identity:  identity,
		peerHost:  peerHost,
		peerstore: peerstore,

		ListenersLocal: newListenersLocal(),
		ListenersP2P:   newListenersP2P(*peerHost),

		Streams: &StreamRegistry{
			Streams:     map[uint64]*Stream{},
			ConnManager: (*peerHost).ConnManager(),
			conns:       map[peer.ID]int{},
		},
	}
}

// CheckProtoExists checks whether a proto handler is registered to
// mux handler
func (p2p *P2P) CheckProtoExists(proto string) bool {
	protos := (*p2p.peerHost).Mux().Protocols()

	for _, p := range protos {
		if p != proto {
			continue
		}
		return true
	}
	return false
}

//Host Return a pointer to peerHost
func (p2p *P2P) Host() *p2phost.Host {
	return p2p.peerHost
}
