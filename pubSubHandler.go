package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gogo/protobuf/proto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func subHandler(ctx context.Context, sub *pubsub.Subscription) {
	log.Info("-- Sub start listen --")
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
		//log.Infof("-->>SUB Recieve: %v type:%s, ID:%s ,Addr: %s ,Extra: %s\n", req.Command,
		//	string(req.Parameters[0]), string(req.Parameters[1]), string(req.Parameters[2]), string(req.Parameters[3]))
	}
}
