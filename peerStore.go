package main

import (
	"errors"

	p2ppeers "github.com/libp2p/go-libp2p-core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

const (
	maxRemotePeers = 3
)

// AddrStringToAddrInfo convert string address to peer AddrInfo
func AddrStringToAddrInfo(addr string) (*p2ppeers.AddrInfo, error) {
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// Extract the peer ID from the multiaddr.
	info, err := p2ppeers.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if nil == info {
		return nil, errors.New("AddrInfo is nil")
	}
	//p.PeersRemote.AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	return info, nil
}
