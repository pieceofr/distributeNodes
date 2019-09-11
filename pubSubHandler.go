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

		if req.Command == cmdPeer {
			var info NodeInfoMessage
			proto.Unmarshal(req.Parameters[0], &info)
			em.Emit(info)
			if err != nil {
				log.Error(err.Error())
			}
			log.Debugf(" --><--  Emit PeerInfo:%v\n", info.ID)
		} else {
			log.Info("unsupported command")
		}
	}
}
