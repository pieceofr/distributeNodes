package main

import (
	libp2pcore "github.com/libp2p/go-libp2p-core"
	libp2prpc "github.com/libp2p/go-libp2p-gorpc"
)

func registerRPC(host *libp2pcore.Host) (*libp2prpc.Server, error) {
	s := libp2prpc.NewServer(*host, "rpc")
	var status NodeStatisticMessage
	err := s.Register(status)
	if err == nil {
		log.Error("expected an error")
		return nil, err
	}
	return s, nil
}
