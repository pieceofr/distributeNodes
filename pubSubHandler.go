package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gogo/protobuf/proto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func (p *PeerNode) subHandler(ctx context.Context, sub *pubsub.Subscription) {
	log.Info("-- Sub start listen --")
	em, err := p.Host.EventBus().Emitter(new(NodeInfoMessage))
	if err != nil {
		panic(err)
	}
	defer em.Close()

	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		req := &NodeMessage{}
		err = proto.Unmarshal(msg.Data, req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		log.Infof("-->>SUB Recieve: %v  ID:%s \n", req.Command, shortID(string(req.Parameters[1])))

		if req.Command == cmdPeer {
			peerInfo := NodeInfoMessage{
				NodeType: NodeType(1),
				ID:       string(req.Parameters[1]),
				Address:  append([]string{}, string(req.Parameters[2])),
				Extra:    string(req.Parameters[3]),
			}
			em.Emit(peerInfo)
			if err != nil {
				log.Error(err.Error())
			}
			log.Debugf(" --><--  Emit PeerInfo:%v\n", peerInfo)
		}
	}
}
