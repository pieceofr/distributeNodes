package main

import (
	"time"

	p2pnet "github.com/libp2p/go-libp2p-core/network"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func (p *PeerNode) networkMonitor() {
	p.Host.Network().Notify(&p2pnet.NotifyBundle{
		ListenF: func(net p2pnet.Network, addr multiaddr.Multiaddr) {
			log.Infof("@@Host: %v is listen at %v\n", addr.String(), time.Now())
		},
		ConnectedF: func(net p2pnet.Network, conn p2pnet.Conn) {
			peerAddr := conn.RemoteMultiaddr()
			log.Infof("@@: Conn: %v Connected at %v\n", peerAddr.String(), time.Now())
		},
		DisconnectedF: func(net p2pnet.Network, conn p2pnet.Conn) {
			peerAddr := conn.RemoteMultiaddr()
			log.Infof("@@Conn: %v Disconnected at %v\n", peerAddr.String(), time.Now())
		},
		OpenedStreamF: func(net p2pnet.Network, stream p2pnet.Stream) {
			peerAddr := stream.Conn().RemoteMultiaddr()
			log.Infof("@@Stream : %v-%v is Opened at %v\n", peerAddr.String(), stream.Protocol(), time.Now())
		},
		ClosedStreamF: func(net p2pnet.Network, stream p2pnet.Stream) {
			peerAddr := stream.Conn().RemoteMultiaddr()
			log.Infof("@@Stream :%v-%v is Closed at %v\n", peerAddr.String(), stream.Protocol(), time.Now())
		},
	})
}
