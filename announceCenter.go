package main

import (
	"errors"
	"fmt"
	"time"

	proto "github.com/golang/protobuf/proto"
)

var broadcastInterval = 15 * time.Second

const (
	cmdPeer = "peer"
)

func (p *PeerNode) announceCenter(shutdown chan struct{}) {
	selfMulticastTimer := time.After(broadcastInterval)
	for {
		select {
		case <-shutdown:
			break
		case <-selfMulticastTimer:
			selfMulticastTimer = time.After(cycleInterval)
			go p.sendSelfNodeMessage()
		}
	}
}

func (p *PeerNode) sendSelfNodeMessage() error {
	// Broadcast self
	if Client == p.NodeInfo.NodeType {
		return nil
	}
	p.NodeInfo.Extra = fmt.Sprintf("%v", time.Now())
	messages := p.makeNodeMessage(&p.NodeInfo)
	for _, msg := range messages {
		err := p.BroadcastStream.Publish(pubsubTopic, msg)
		if err != nil {
			log.Errorf("broadcast error: %v\n", err)
			return err
		}
		log.Infof("<<--- broadcasting SELF : %v\n", shortID(p.NodeInfo.ID))
	}

	return nil
}

func (p *PeerNode) sendPeerNodeMessage(peersInfo map[string]NodeInfoMessage) error {
	for key, node := range peersInfo {
		p.Mutex.Lock()
		messages := p.makeNodeMessage(&node)
		if nil == messages {
			return errors.New("make proto message error")
		}
		delete(peersInfo, key)
		p.Mutex.Unlock()
		if node.ID != p.NodeInfo.ID { // Node is not self
			for _, msg := range messages {
				err := p.BroadcastStream.Publish(pubsubTopic, msg)
				if err != nil {
					log.Warnf("broadcast error: %v\n", err)
					return err
				}
				log.Infof("<<--- broadcasting PEER : %v\n", node)
			}
		}
	}
	return nil
}

func (p *PeerNode) makeNodeMessage(node *NodeInfoMessage) [][]byte {
	var params [][]byte
	var messages [][]byte
	for _, addr := range node.Address {
		params = append(params, []byte(fmt.Sprintf("%v", node.NodeType)), []byte(node.ID), []byte(addr), []byte(node.Extra))

		req := &NodeMessage{
			Command:    cmdPeer,
			Parameters: params,
		}
		msgBytes, err := proto.Marshal(req)
		if err != nil {
			log.Warnf("broadcast proto error: %v\n", err)
			return nil
		}
		messages = append(messages, msgBytes)
	}
	return messages
}
