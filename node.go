package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	p2pnet "github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tls "github.com/libp2p/go-libp2p-tls"
	ma "github.com/multiformats/go-multiaddr"
	multiaddr "github.com/multiformats/go-multiaddr"
	//manet "github.com/multiformats/go-multiaddr-net"
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
const pubsubTopic = "/peer/announce/1.0.0"
const nodeProtocol = "/p2p"

var (
	privateKeyFileName  = "peer.prv"
	discoveryRetryTimes = 5
)

// PeerNode peer Node struct
type PeerNode struct {
	NodeType
	PublicIP []string
	Port     string
	Identity
	Host            libp2pcore.Host
	Streams         []p2pnet.Stream
	PeersRemote     peerstore.Peerstore
	NodeInfo        NodeInfoMessage
	Shutdown        chan struct{}
	Mutex           *sync.Mutex
	Handlers        []NodeStreamHandler
	BroadcastStream *pubsub.PubSub
}

//Init/setup a peer node
func (p *PeerNode) setup(cfg config) error {
	p.Mutex = &sync.Mutex{}
	p.NodeType = NodeType(cfg.NodeType)
	p.PublicIP = cfg.PublicIP
	if cfg.StaticIdentity.UseStatic { // Use static identity from file
		log.Info("Use Static Identity")
		privateKeyFileName = cfg.StaticIdentity.KeyFile
		if loadErr := p.LoadIdentity(privateKeyFileName); loadErr != nil {
			log.Error(loadErr.Error())
			return loadErr
		}
		log.Info("Load Static Identity")
	} else { // Random Generate Identity
		if genRandErr := p.NewRandomNode(); genRandErr != nil {
			log.Error(genRandErr.Error())
			return genRandErr
		}
	}

	p.Shutdown = make(chan struct{})
	p.Port = strconv.Itoa(cfg.Port)
	if Servant == p.NodeType || Server == p.NodeType {
		if createHostErr := p.NewServantHost(); createHostErr != nil {
			log.Error("createHostErr")
			panic(createHostErr)
		}
	} else {
		newHost, err := libp2p.New(
			context.Background(),
			libp2p.Identity(p.Identity.PrvKey),
			libp2p.Security(tls.ID, tls.New),
		)
		if err != nil {
			return err
		}
		p.Port = "" // If user the port ie. 2136. this will be check with owns port
		p.Host = newHost
	}
	p.networkMonitor()
	ps, err := pubsub.NewGossipSub(context.Background(), p.Host)
	if err != nil {
		panic(err)
	}
	p.BroadcastStream = ps
	go p.peersTable(p.Shutdown)

	maAddrs := []ma.Multiaddr{}
	for _, ip := range p.PublicIP {
		addr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%v%s/%s", ip, p.Port, nodeProtocol, p.Host.ID().Pretty()))
		if nil == err {
			maAddrs = append(maAddrs, addr)
		}
	}

	p.NodeInfo = NodeInfoMessage{
		NodeType: uint32(p.NodeType),
		ID:       p.Host.ID().String(),
		Addrs:    &Addrs{Address: GetBytesFromMultiaddr(maAddrs)},
	}
	log.Info("Exist the Initialization")
	return nil
}

func (p *PeerNode) run() {
	//go p.BusReciever(p.Shutdown)
	go registerRPC(&p.Host)
	sub, err := p.BroadcastStream.Subscribe(pubsubTopic)
	go p.subHandler(context.Background(), sub)
	if err != nil {
		panic(err)
	}
	go p.announceCenter(p.Shutdown)
	if Servant == p.NodeType || Server == p.NodeType {
		log.Info("Servant go listen routine")
		go p.Listen()
	}
	go p.ConnectToFixPeer()
}

//SaveIdentity save private key to file
func (p *PeerNode) SaveIdentity(file string) error {
	if "" == file {
		file = privateKeyFileName
	}
	keyStr, err := p.Identity.MarshalPrvKey()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(file, []byte(keyStr), 0600); err != nil {
		return err
	}
	return nil
}

//LoadIdentity Load private key from file
func (p *PeerNode) LoadIdentity(file string) error {
	if len(file) == 0 {
		file = privateKeyFileName
	}
	keyBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = p.UnmarshalPrvKey(keyBytes)
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

// NewServantHost Create a NewHost of server
func (p *PeerNode) NewServantHost() error {
	listenAddrIPV4, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", "0.0.0.0", p.Port))
	if err != nil {
		panic(err)
	}

	listenAddrIPV6, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip6/%s/tcp/%s", "::", p.Port))
	if err != nil {
		panic(err)
	}

	//listenAddrs := listenAddrIPV4.Encapsulate(listenAddrIPV6)
	newHost, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(listenAddrIPV4, listenAddrIPV6),
		libp2p.Identity(p.Identity.PrvKey),
		libp2p.Peerstore(p.PeersRemote),
		libp2p.Security(tls.ID, tls.New),
	)
	if err != nil {
		return err
	}
	for _, a := range newHost.Addrs() {
		log.Infof("New Host Address: %s/%v/%s\n", a, "p2p", newHost.ID())
	}

	p.Host = newHost
	return nil
}

//Listen host start to listen
func (p *PeerNode) Listen() error {
	if p.Host == nil {
		return errors.New("NoHost")
	}
	var handleStream NodeStreamHandler
	handleStream.HandlerType = ListenerHander
	p.Mutex.Lock()
	p.Handlers = append(p.Handlers, handleStream)
	handleStream.Setup(len(p.Handlers), p.NodeInfo, &p.Host)
	p.Mutex.Unlock()
	p.Host.SetStreamHandler(nodeProtocol, handleStream.Handler)

	for _, la := range p.Host.Network().ListenAddresses() {
		if _, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			return err
		}
		log.Infof("LISTEN ON ===")
		for _, ip := range p.PublicIP {
			log.Infof("/ip4/%s/tcp/%v/%v/%s' on another console.\n", ip, p.Port, "p2p", p.Host.ID().Pretty())
		}
		log.Infof("%s\n", la.String())
		log.Info("\nWaiting for incoming connection\n\n")
	}
	<-make(chan struct{})
	return nil
}

//ConnectToFixPeer connect to fixed peer
func (p *PeerNode) ConnectToFixPeer() error {
	for i := 0; i < discoveryRetryTimes; i++ {
		addr, err := GetAServer()
		if err != nil {
			return err
		}
		maAddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil
		}
		err = p.ConnectTo(maAddr)
		if nil == err {
			i = 10
			log.Infof("HAS CONNECT ED TO : ", addr)
		}
		log.Error(err.Error())
		time.Sleep(3 * time.Second)
	}
	return nil
}

//ConnectTo connect to servant or server
func (p *PeerNode) ConnectTo(address ma.Multiaddr) error {
	if p.NodeType == Server {
		return errors.New("NodeType:Server.Do not have ability to connect")
	}
	if p.IsSameMa(address) {
		return errSameNode
	}

	info, err := addrMaToAddrInfo(address)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	p.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	//	p.PeersRemote.AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	log.Debugf("Connecting .... ID:%s, address:%s TTL:%d", info.ID.String(), info.Addrs[0].String(), peerstore.PermanentAddrTTL)

	s, err := p.Host.NewStream(context.Background(), info.ID, nodeProtocol)
	if err != nil {
		return err
	}
	p.Streams = append(p.Streams, s)
	log.Infof("NEW STREAM ID:%s, address:%s TTL:%d", info.ID.String(), info.Addrs[0].String(), peerstore.PermanentAddrTTL)
	var handleStream NodeStreamHandler
	handleStream.HandlerType = ClientHandler
	p.Mutex.Lock()
	p.Handlers = append(p.Handlers, handleStream)
	handleStream.Setup(len(p.Handlers), p.NodeInfo, &p.Host)
	p.Mutex.Unlock()
	handleStream.Handler(s)
	<-make(chan struct{})
	return nil
}

//Reset is to close a peer node
func (p *PeerNode) Reset() {
	p.NodeType = Servant
	p.PublicIP = []string{}
	p.Port = ""
	p.Identity.PrvKey = nil

	if p.Host != nil {
		p.Host.Close()
	}
	log.Warn("Reset Peer Node")
}
