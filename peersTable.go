package main

import (
	"time"

	p2ppeers "github.com/libp2p/go-libp2p-core/peer"
)

const (
	numberOfConnections = 3
	printInterval       = 40 * time.Second
)

var messageQ map[string]NodeInfoMessage

func (p *PeerNode) peersTable(shutdown chan struct{}) {
	eventbus := p.Host.EventBus()
	messageQ = make(map[string]NodeInfoMessage, 100)
	sub, err := eventbus.Subscribe(new(NodeInfoMessage))
	if err != nil {
		panic("Can not subscribe NodeInfoMessage")
	}
	selfPrintTimer := time.After(printInterval)
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
		case <-selfPrintTimer:
			selfPrintTimer = time.After(broadcastInterval)
			p.Mutex.Lock()
			log.Infof("----------peerTable List Start-----------\n")
			for key, node := range messageQ {
				log.Infof("[---] ID:%v  NODE:%v\n", key, node)
			}
			log.Infof("----------peerTable List End-----------\n")
			p.Mutex.Unlock()
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
				for _, addr := range GetMultiAddrsFromBytes(node.Addrs.Address) {
					go p.ConnectTo(addr)
					log.Infof("@@Try to Connect to %s\n", addr.String())
				}
			}
		}
		count++
	}
	return nil
}
