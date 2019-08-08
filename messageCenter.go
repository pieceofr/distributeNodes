package main

import (
	"errors"
	"fmt"
	"time"

	proto "github.com/golang/protobuf/proto"
)

var broadcastInterval = 10 * time.Second

func (p *PeerNode) announceCenter(shutdown chan struct{}) {
	eventbus := p.Host.EventBus()
	sub, err := eventbus.Subscribe(new(NodeInfoMessage))
	defer sub.Close()
	if err != nil {
		panic("Can not subscribe NodeInfoMessage")
	}
	messageQ := make(map[string]NodeInfoMessage, 100)
	cycleTimer := time.After(broadcastInterval)
	for {
		select {
		case <-shutdown:
			break
		case e := <-sub.Out():
			switch e.(type) {
			case NodeInfoMessage:
				log.Debugf("[announce]recieve:%v\n", (e.(NodeInfoMessage)))
				messageQ[e.(NodeInfoMessage).ID] = e.(NodeInfoMessage)
			}
		case <-cycleTimer:
			cycleTimer = time.After(cycleInterval)
			go p.sendSelfNodeMessage()
			//Help to Broadcast Peers
			log.Info("***Help peers broadcast")
			go p.sendPeerNodeMessage(messageQ)
		}
	}
}

func (p *PeerNode) sendSelfNodeMessage() error {
	// Broadcast self
	if Client == p.NodeInfo.NodeType {
		return nil
	}
	p.NodeInfo.Extra = fmt.Sprintf("%v", time.Now())
	msgBytes := p.makeNodeMessage(&p.NodeInfo)
	err := p.BroadcastStream.Publish(pubsubTopic, msgBytes)
	if err != nil {
		log.Errorf("broadcast error: %v\n", err)
		return err
	}
	log.Infof("<<--- broadcasting SELF : %v\n", shortID(p.NodeInfo.ID))
	return nil
}

func (p *PeerNode) sendPeerNodeMessage(peersInfo map[string]NodeInfoMessage) error {
	for key, node := range peersInfo {
		p.Mutex.Lock()
		msgBytes := p.makeNodeMessage(&node)
		if nil == msgBytes {
			return errors.New("make proto message error")
		}
		delete(peersInfo, key)
		p.Mutex.Unlock()
		if node.ID != p.NodeInfo.ID {
			err := p.BroadcastStream.Publish(pubsubTopic, msgBytes)
			if err != nil {
				log.Warnf("broadcast error: %v\n", err)
				return err
			}
			log.Infof("<<--- broadcasting PEER : %v\n", node)
		}
	}
	return nil
}

func (p *PeerNode) makeNodeMessage(node *NodeInfoMessage) []byte {
	var params [][]byte
	params = append(params, []byte(fmt.Sprintf("%v", node.NodeType)), []byte(node.ID), []byte(node.Address), []byte(node.Extra))
	req := &NodeMessage{
		Command:    "peer",
		Parameters: params,
	}
	msgBytes, err := proto.Marshal(req)
	if err != nil {
		log.Warnf("broadcast proto error: %v\n", err)
		return nil
	}
	return msgBytes
}
