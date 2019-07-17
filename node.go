package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/libp2p/go-libp2p"
	libp2pcore "github.com/libp2p/go-libp2p-core"
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

	if cfg.StaticIdentity.UseStatic {
		log.Info("Use Static Identity")
		p.NodeType = NodeType(cfg.NodeType)
		p.Port = strconv.Itoa(cfg.Port)
		privateKeyFileName = cfg.StaticIdentity.KeyFile
		if loadErr := p.LoadIdentity(); loadErr != nil {
			return loadErr
		}
	} else {
		p.NodeType = NodeType(cfg.NodeType)
		if genRandErr := p.NewRandomNode(); genRandErr != nil {
			return genRandErr
		}
	}
	if Servant == p.NodeType {
		if createHostErr := p.NewHost(); createHostErr != nil {
			return createHostErr
		}
	}
	log.Info(p.Port)

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
	if p.Host != nil {
		return errors.New("NoHost")
	}
	p.Host.SetStreamHandler("/peer/1.0.0", handleStream)
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
