package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
)

//NodeType present the type of a node
type NodeType int

const (
	//Servant acts as both server and client
	Servant NodeType = iota
	//Client acts as a client only
	Client
	//Server acts as a server only
	Server
)

var (
	privateKeyFileName = "peer.prv"
)

// PeerNode peer Node struct
type PeerNode struct {
	NodeType
	PublicIP string
	Port     string
	Identity
	Host libp2pcore.Host
}

// Connection is a p2p connection struct
/*
type Connection struct {
	privateKey crypto.PrivKey
	connHost   libp2pcore.Host
}*/

//Init a peer node
func (p *PeerNode) Init(cfg config) error {
	p.NodeType = NodeType(cfg.NodeType)
	p.Port = strconv.Itoa(cfg.Port)
	p.PublicIP = cfg.PublicIP
	if cfg.StaticIdentity.UseStatic {
		log.Info("Use Static Identity")
		privateKeyFileName = cfg.StaticIdentity.KeyFile
		if loadErr := p.LoadIdentity(); loadErr != nil {
			return loadErr
		}
	} else {
		if genRandErr := p.NewRandomNode(); genRandErr != nil {
			return genRandErr
		}
	}
	if Servant == p.NodeType || Server == p.NodeType {
		if createHostErr := p.NewHost(); createHostErr != nil {
			log.Error("createHostErr")
			return createHostErr
		}
		err := p.Listen()
		if err != nil {
			log.Error(err.Error())
			return err
		}
	} else {
		newHost, err := libp2p.New(
			context.Background(),
			libp2p.Identity(p.Identity.PrvKey),
		)

		if err != nil {
			return err
		}
		p.Host = newHost
	}
	go p.ConnectToFixPeer()
	log.Info("Exist the Initialization")
	return nil
}

//SaveIdentity save private key to file
func (p *PeerNode) SaveIdentity() error {
	keyStr, err := p.Identity.MarshalPrvKey()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(privateKeyFileName, []byte(keyStr), 0600); err != nil {
		return err
	}
	return nil
}

//LoadIdentity Load private key from file
func (p *PeerNode) LoadIdentity() error {
	keyBytes, err := ioutil.ReadFile(privateKeyFileName)
	if err != nil {
		return err
	}
	err = p.UnmarshalPrvKey(string(keyBytes))
	if err != nil {
		return err
	}
	return nil
}

//NewRandomNode generate a default Connection
func (p *PeerNode) NewRandomNode() error {
	// generate a random Identity
	var nodeID Identity
	err := nodeID.randIdentity()
	if err != nil {
		return nil
	}
	p.Identity = nodeID
	port, err := randomPort()
	if err != nil {
		return err
	}
	p.Port = strconv.Itoa(port)
	return nil
}

//NewHost Create a NewHost of server
func (p *PeerNode) NewHost() error {
	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", "0.0.0.0", p.Port))
	if err != nil {
		return err
	}
	newHost, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(p.Identity.PrvKey),
	)
	if err != nil {
		return err
	}
	p.Host = newHost
	return nil
}

//Listen host start to listen
func (p *PeerNode) Listen() error {
	if p.Host == nil {
		return errors.New("NoHost")
	}
	p.Host.SetStreamHandler("/p2p/1.0.0", handleStream)
	for _, la := range p.Host.Network().ListenAddresses() {
		if _, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			return err
		}
		log.Infof("Listen port:%s", p.Port)
		log.Infof("Run './chat -d /ip4/%s/tcp/%v/p2p/%s' on another console.\n", p.PublicIP, p.Port, p.Host.ID().Pretty())
		log.Info("\nWaiting for incoming connection\n\n")
	}
	<-make(chan struct{})
	return nil
}

//ConnectToFixPeer connect to fixed peer
func (p *PeerNode) ConnectToFixPeer() error {
	addr, err := GetAServer()
	if err != nil {
		return err
	}
	log.Infof("Client is connect to %s", addr)
	for i := 0; i < 5; i++ {
		err := p.ConnectTo(addr)
		if nil == err {
			i = 10
		}
		log.Errorf("The Error is %s", err.Error())
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Info("Exist ConnectToFixPeer")
		return err
	}

	return nil
}

//ConnectTo connect to servant or server
func (p *PeerNode) ConnectTo(address string) error {
	if p.NodeType == Server {
		return errors.New("NodeType:Server.Do not have ability to connect")
	}
	if p.IsSameNode(address) {
		return ErrSameNode
	}

	maddr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if info == nil {
		return errors.New("info is nil")
	}
	log.Debugf("ID:%s, address:%s TTL:%d", info.ID.String(), info.Addrs[0].String(), peerstore.PermanentAddrTTL)

	time.Sleep(1 * time.Second)
	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	p.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := p.Host.NewStream(context.Background(), info.ID, "/p2p/1.0.0")
	if err != nil {
		panic(err)
	}

	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// Create a thread to read and write data.
	go writeData(rw)
	go readData(rw)

	<-make(chan struct{})
	return nil
}

//Reset is to close a peer node
func (p *PeerNode) Reset() {
	p.NodeType = Servant
	p.PublicIP = ""
	p.Port = ""
	p.Identity.PrvKey = nil
	if p.Host != nil {
		p.Host.Close()
	}
	log.Warn("Reset Peer Node")
}

//IsSameNode check if provided address has the same ip and port with peerNode
func (p *PeerNode) IsSameNode(addr string) bool {
	elems := strings.Split(addr, "/")
	if elems[2] != "" && elems[2] == p.PublicIP {
		if elems[4] == p.Port {
			log.Debugf("addr[2]:%s  p.Public:%s   elemes[4]:%s p.Port:%s", elems[2], p.PublicIP, elems[4], p.Port)
			return true
		}
	}
	return false
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
