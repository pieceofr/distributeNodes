package main

import (
	"context"
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	libp2pcore "github.com/libp2p/go-libp2p-core"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// PeerNode peer Node struct
type PeerNode struct {
	NodeType int //0:servant 1:client 2:server
	PubKey   []byte
	Host     libp2pcore.Host
	Port     int
}

// Connection is a p2p connection struct
type Connection struct {
	privateKey crypto.PrivKey
	connHost   libp2pcore.Host
}

//NewRandomNode generate a default Connection
func (c *Connection) NewRandomNode(nodeType int) error {
	// generate a random Identity
	var thisNode Identity
	thisNode.randIdentity()
	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", KnownServant[0].Host, KnownServant[0].Port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	newHost, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(thisNode.PrvKey),
	)
	if err != nil {
		return err
	}
	c.connHost = newHost
	return nil
}
