package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	ma "github.com/multiformats/go-multiaddr"
	multiaddr "github.com/multiformats/go-multiaddr"
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

//IsSameMa Check if PeerNode has the same address
func (p *PeerNode) IsSameMa(addr ma.Multiaddr) bool {
	for _, a := range p.Host.Addrs() {
		if strings.Contains(addr.String(), a.String()) {
			log.Infof("The same string listenAddr:%s NodeInfoMessage:%s", a.String(), addr.String())
			return true
		}
	}
	return false
}

//IsSameNode check if provided address has the same ip and port with peerNode
func (p *PeerNode) IsSameNode(addr string) bool {
	elems := strings.Split(addr, "/")
	if elems[2] != "" {
		for _, ip := range p.PublicIP {
			if elems[2] == ip {
				if elems[4] == p.Port {
					log.Debugf("addr[2]:%s  p.Public:%s   elemes[4]:%s p.Port:%s", elems[2], ip, elems[4], p.Port)
					return true
				}
			}
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

func randNum(num int) int {
	if 0 == num {
		return -1
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}

func pickNumbersInSize(size, num int) []bool {
	if size < num || size <= 0 || num <= 0 {
		log.Warn("invalid input or zero numer")
		return []bool{}
	}
	pickupTable := make([]bool, size)
	for i := 0; i < num; {
		if pickNum := randNum(size); -1 == pickNum {
			return []bool{}
		} else {
			if false == pickupTable[pickNum] {
				pickupTable[pickNum] = true
				i++
			}
		}
	}
	ret := ""
	for _, val := range pickupTable {
		ret = fmt.Sprintf("%s %v", ret, val)
	}
	return pickupTable
}

// GetMultiAddrsFromBytes take  [][]byte listeners and convert them into []Multiaddr format
func GetMultiAddrsFromBytes(listners [][]byte) []ma.Multiaddr {
	var maAddrs []ma.Multiaddr
	for _, addr := range listners {
		maAddr, err := ma.NewMultiaddrBytes(addr)
		if nil == err {
			maAddrs = append(maAddrs, maAddr)
		}
	}
	return maAddrs
}

// GetBytesFromMultiaddr take []Multiaddr format listeners and convert them into   [][]byte
func GetBytesFromMultiaddr(listners []ma.Multiaddr) [][]byte {
	var byteAddrs [][]byte
	for _, addr := range listners {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}
	return byteAddrs
}
