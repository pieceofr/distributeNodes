package main

import (
	"context"
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
)

func newHost(port string) (h1 host.Host) {
	h1, _ = libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", port)),
	)
	return
}

//RPCTest For Test purpose
type RPCTest struct {
}
type RPCArg struct {
	Send string
}

//RPCPingPong tell user communication is works
func (m *RPCTest) RPCPingPong(ctx context.Context, recv *string, resp *string) error {
	pong := "pong"
	if *recv != "ping" {
		pong = "What do you say ?"
	}
	resp = &pong
	return nil
}

/*
func TestRpcRegister(t *testing.T) {
	h1 := newHost("12137")
	defer h1.Close()
	s := libp2prpc.NewServer(h1, "rpc")
	var test RPCTest

	err := s.Register(&test)
	assert.NoError(t, err, "RPC Register fail")
}


func TestRpcConnection(t *testing.T) {
	servHost := newHost("12137")
	s := libp2prpc.NewServer(servHost, "rpc")
	var test RPCTest
	err := s.Register(&test)
	assert.NoError(t, err, "RPC Register fail")
	clientHost := newHost("12138")
	rpcClient := libp2prpc.NewClient(clientHost, "rpc")
	var resp string
	RPCArg := RPCArg{Send: "ping"}
	err = rpcClient.Call(servHost.ID(), "RPCTest", "RPCPingPong", &RPCArg, &resp)
	assert.NoError(t, err, "RPC Call  fail")

}
*/
