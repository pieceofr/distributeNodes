package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	libp2pcore "github.com/libp2p/go-libp2p-core"

	"github.com/libp2p/go-libp2p-core/network"
	libp2pnetwork "github.com/libp2p/go-libp2p-core/network"
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
	privateKeyFileName  = "peer.prv"
	discoveryRetryTimes = 5
	broadcastInterval   = 3 * time.Second
)

// PeerNode peer Node struct
type PeerNode struct {
	NodeType
	PublicIP string
	Port     string
	Identity
	Host          libp2pcore.Host
	Streams       []libp2pnetwork.Stream
	PeersRemote   peerstore.Peerstore
	PeersListener peerstore.Peerstore
	NodeInfo      NodeInfoMessage
	Shutdown      chan<- struct{}
	Mutex         *sync.Mutex
	Handlers      []NodeStreamHandler
}

//Init a peer node
func (p *PeerNode) Init(cfg config) error {
	p.Mutex = &sync.Mutex{}
	p.NodeType = NodeType(cfg.NodeType)
	p.PublicIP = cfg.PublicIP
	if cfg.StaticIdentity.UseStatic {
		log.Info("Use Static Identity")
		privateKeyFileName = cfg.StaticIdentity.KeyFile
		if loadErr := p.LoadIdentity(privateKeyFileName); loadErr != nil {
			return loadErr
		}
		log.Info("Load Static Identity")
	} else {
		if genRandErr := p.NewRandomNode(); genRandErr != nil {
			log.Error(genRandErr.Error())
			return genRandErr
		}
	}
	shutdown := make(chan struct{})
	p.Shutdown = shutdown
	go p.BusReciever(shutdown)
	p.Port = strconv.Itoa(cfg.Port)
	if Servant == p.NodeType || Server == p.NodeType {
		if createHostErr := p.NewHost(); createHostErr != nil {
			log.Error("createHostErr")
			panic(createHostErr)
		}
	} else {
		newHost, err := libp2p.New(
			context.Background(),
			libp2p.Identity(p.Identity.PrvKey),
			libp2p.NATPortMap(),
		)
		if err != nil {
			return err
		}
		p.Port = "" // If user the port ie. 2136. this will be check with owns port
		p.Host = newHost
	}
	for _, addr := range p.Host.Addrs() {
		log.Infof("Host Address: \n", addr.String())
	}

	p.NodeInfo = NodeInfoMessage{NodeType: p.NodeType, ID: fmt.Sprintf("%v", p.Host.ID()), Address: fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", p.PublicIP, p.Port, p.Host.ID().Pretty())}
	log.Infof("NodeAddress:%s\n", p.NodeInfo.Address)

	if Servant == p.NodeType || Server == p.NodeType {
		log.Info("Servant go listen routine")
		go p.Listen()

	}
	go p.ConnectToFixPeer()

	log.Info("Exist the Initialization")
	return nil
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
	pkey, err := p.MarshalPublicKey()
	log.Debugf("NewRandomNode:PublicKey:%s   Port:%s", pkey, p.Port)
	return nil
}

//NewHost Create a NewHost of server
func (p *PeerNode) NewHost() error {
	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", "0.0.0.0", p.Port))
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	newHost, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(p.Identity.PrvKey),
		libp2p.Peerstore(p.PeersListener),
		libp2p.NATPortMap(),
	)
	if err != nil {
		return err
	}
	for _, a := range newHost.Addrs() {
		log.Infof("New Host Address: %s/%v/%s\n", a, "p2p", peer.IDB58Encode(newHost.ID()))
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
	handleStream.Setup(len(p.Handlers), p.NodeInfo)
	p.Mutex.Unlock()
	p.Host.SetStreamHandler("/p2p/1.0.0", handleStream.Handler)

	for _, la := range p.Host.Network().ListenAddresses() {
		if _, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			return err
		}
		log.Infof("LISTEN ON ===")
		log.Infof("/ip4/%s/tcp/%v/%v/%s' on another console.\n", p.PublicIP, p.Port, "p2p", p.Host.ID().Pretty())
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
		err = p.ConnectTo(addr)
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
func (p *PeerNode) ConnectTo(address string) error {
	if p.NodeType == Server {
		return errors.New("NodeType:Server.Do not have ability to connect")
	}
	if p.IsSameNode(address) {
		return errSameNode
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

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.

	shortAddr, err := multiaddr.NewMultiaddr(addrToConnAddr(address))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if p.IsPeerExisted(shortAddr) {
		if p.Host.Network().Connectedness(info.ID) != network.Connected {
			connectErr := p.Host.Connect(context.Background(), *info)
			if connectErr != nil {
				log.Errorf("RECONNECT Stream Error:%v", ErrCombind(errReconnectStream, connectErr))
				return ErrCombind(errReconnectStream, connectErr)
			} else {
				log.Infof("ID %s  is not CONNECTED!\n", shortID(info.ID.String()))
			}
		}
		return errPeerHasInPeerStore
	}

	p.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	//	p.PeersRemote.AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	log.Debugf("Connecting .... ID:%s, address:%s TTL:%d", info.ID.String(), info.Addrs[0].String(), peerstore.PermanentAddrTTL)

	s, err := p.Host.NewStream(context.Background(), info.ID, "/p2p/1.0.0")
	if err != nil {
		return err
	}
	p.Streams = append(p.Streams, s)
	log.Infof("NEW STREAM ID:%s, address:%s TTL:%d", info.ID.String(), info.Addrs[0].String(), peerstore.PermanentAddrTTL)
	var handleStream NodeStreamHandler
	handleStream.HandlerType = ClientHandler
	p.Mutex.Lock()
	p.Handlers = append(p.Handlers, handleStream)
	handleStream.Setup(len(p.Handlers), p.NodeInfo)
	p.Mutex.Unlock()
	handleStream.Handler(s)
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

//BusReciever recieve message from other component
func (p *PeerNode) BusReciever(shutdown <-chan struct{}) {
	queue := Bus.TestQueue.Chan()
	cycleTimer := time.After(broadcastInterval)
	messageQ := make(map[string]NodeInfoMessage, 100)
	for {
		select {
		case <-shutdown:
			break
		case <-cycleTimer:
			cycleTimer = time.After(cycleInterval)
			for _, val := range messageQ {
				if !p.IsSameNode(val.Address) && p.NodeType != Client {
					log.Info("--->Help Peers Broadcasting ")
					val.Extra = fmt.Sprintf("%v", time.Now())
					Bus.Broadcast.Send("peer", []byte(fmt.Sprintf("%v", val.NodeType)), []byte(val.ID), []byte(val.Address), []byte(val.Extra))
				}
			}
			break
		case item := <-queue:
			if item.Command != "peer" {
				log.Error("Command is not Peer")
				break
			}
			var peerInfo NodeInfoMessage
			nType, err := strconv.Atoi(string(item.Parameters[0]))
			if err != nil {
				log.Error(err.Error())
				continue
			}
			peerInfo.NodeType = NodeType(nType)
			peerInfo.ID = string(item.Parameters[1])
			peerInfo.Address = string(item.Parameters[2])
			peerInfo.Extra = string(item.Parameters[3])

			log.Infof("from queue: %q  %s\n", item.Command, peerInfo.Address)
			if peerInfo.NodeType != Client {
				go p.ConnectTo(peerInfo.Address)
				messageQ[peerInfo.ID] = peerInfo
				//Bus.Broadcast.Send("testing", []byte(fmt.Sprintf("%v", time.Now())))
			}

		}
	}
}
