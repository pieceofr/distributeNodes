package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/multiformats/go-multiaddr"
)

//IsPeerExisted peer is existed in the Peerstore
func (p *PeerNode) IsPeerExisted(newAddr multiaddr.Multiaddr) bool {
	for _, ID := range p.Host.Peerstore().Peers() {
		for _, addr := range p.Host.Peerstore().PeerInfo(ID).Addrs {
			//	log.Debugf("peers in PeerStore:%s     NewAddress:%s\n", addr.String(), newAddr.String())
			if addr.Equal(newAddr) {
				log.Info("Peer is in PeerStore")
				return true
			}
		}
	}
	return false
}

//IsSameNode check if provided address has the same ip and port with peerNode
func (p *PeerNode) IsSameNode(addr string) bool {
	elems := strings.Split(addr, "/")
	if elems[2] != "" && elems[2] == p.PublicIP {
		if elems[4] == p.Port {
			//	log.Debugf("addr[2]:%s  p.Public:%s   elemes[4]:%s p.Port:%s", elems[2], p.PublicIP, elems[4], p.Port)
			return true
		}
	}
	return false
}

func (p *PeerNode) getPeerWriter() []*bufio.ReadWriter {
	var writers []*bufio.ReadWriter
	for _, handler := range p.Handlers {
		writers = append(writers, handler.ReadWriter)
	}
	return writers
}

func randomPort() (int, error) {
	// Bitmark Open Port from 12130-12150
	// 12130-12136 reserve for servant node
	// random port open from 12137 - 12150
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return 0, err
	}
	port := 12137 + int(math.Mod(float64(b[0]), float64(13)))
	return port, nil
}
