package main

import (
	"time"

	p2pnet "github.com/libp2p/go-libp2p-core/network"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var connCount Satistic
var streamCount Satistic

func (p *PeerNode) networkMonitor() {
	p.Host.Network().Notify(&p2pnet.NotifyBundle{
		ListenF: func(net p2pnet.Network, addr multiaddr.Multiaddr) {
			log.Debugf("@@Host: %v is listen at %v\n", addr.String(), time.Now())
		},
		ConnectedF: func(net p2pnet.Network, conn p2pnet.Conn) {
			peerAddr := conn.RemoteMultiaddr()
			log.Debugf("@@: Conn: %v Connected at %v\n", peerAddr.String(), time.Now())
			connCount.Add(1)
		},
		DisconnectedF: func(net p2pnet.Network, conn p2pnet.Conn) {
			peerAddr := conn.RemoteMultiaddr()
			log.Infof("@@Conn: %v Disconnected at %v\n", peerAddr.String(), time.Now())
			connCount.Sub(1)
		},
		OpenedStreamF: func(net p2pnet.Network, stream p2pnet.Stream) {
			peerAddr := stream.Conn().RemoteMultiaddr()
			log.Debugf("@@Stream : %v-%v is Opened at %v\n", peerAddr.String(), stream.Protocol(), time.Now())
			streamCount.Add(1)
		},
		ClosedStreamF: func(net p2pnet.Network, stream p2pnet.Stream) {
			peerAddr := stream.Conn().RemoteMultiaddr()
			log.Infof("@@Stream :%v-%v is Closed at %v\n", peerAddr.String(), stream.Protocol(), time.Now())
			streamCount.Sub(1)
		},
	})
}
