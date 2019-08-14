package main

import (
	p2ppeers "github.com/libp2p/go-libp2p-core/peer"
)

const numberOfConnections = 3

var messageQ map[string]NodeInfoMessage

func (p *PeerNode) peersTable(shutdown chan struct{}) {
	eventbus := p.Host.EventBus()
	messageQ = make(map[string]NodeInfoMessage, 100)
	sub, err := eventbus.Subscribe(new(NodeInfoMessage))
	if err != nil {
		panic("Can not subscribe NodeInfoMessage")
	}
	defer sub.Close()
	for {
		select {
		case <-shutdown:
			break
		case e := <-sub.Out():
			switch e.(type) {
			case NodeInfoMessage:
				log.Infof(" --><-- Recieve NodeInfoMessage:%v\n", (e.(NodeInfoMessage)))
				p.Mutex.Lock()
				messageQ[e.(NodeInfoMessage).ID] = e.(NodeInfoMessage)
				p.Mutex.Unlock()
				p.connectPeerCandidates(&messageQ)
			}
		}
	}
}

func (p *PeerNode) connectPeerCandidates(peersTable *map[string]NodeInfoMessage) error {
	pickNumber := numberOfConnections
	if len(*peersTable) < numberOfConnections {
		pickNumber = len(*peersTable)
	}
	pickedNumbers := pickNumbersInSize(len(*peersTable), pickNumber)
	count := 0
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	for key, node := range *peersTable {
		if count < len(pickedNumbers) && pickedNumbers[count] {
			if p2ppeers.ID(key) != p.Host.ID() {
				go p.ConnectTo(node.Address)
				log.Infof("@@Try to Connect to %s\n", node.Address)
			}
		}
		count++
	}
	return nil
}
