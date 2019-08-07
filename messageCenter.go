package main

import (
	"time"
)

var broadcastInterval = 20 * time.Second

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
			//Help to Broadcast Peers
			log.Info("***Help peer broadcast")
		}
	}
}
