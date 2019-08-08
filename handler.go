package main

import (
	"bufio"
	"time"

	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/network"
)

const (
	//MessageSeperator seperator between header and content
	MessageSeperator = "::"
	//HeaderAnnounceSelf s the header to announce self
	HeaderAnnounceSelf = "peer"
	cycleInterval      = 10 * time.Second
)

//HandlerType type of Client:0 Listener:1
type HandlerType int

const (
	//ClientHandler handler
	ClientHandler HandlerType = iota
	//ListenerHander acts as a client only
	ListenerHander
)

var handlerNum int

//NodeStreamHandler for  node to handle stream
type NodeStreamHandler struct {
	Host *libp2pcore.Host
	HandlerType
	Stream     network.Stream
	ReadWriter *bufio.ReadWriter
	NodeInfoMessage
	ID       int
	Shutdown chan<- struct{}
}

//Setup setup Handler
func (h *NodeStreamHandler) Setup(id int, info NodeInfoMessage, host *libp2pcore.Host) {
	h.Host = host
	h.ID = id
	h.NodeInfoMessage = info
	log.Infof("New Stream Handler:%d\n\n", id)
}

// Handler  for streamHandler
func (h *NodeStreamHandler) Handler(s network.Stream) {
	globalMutex.Lock()
	handlerNum++
	h.Stream = s
	h.ReadWriter = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	shutdown := make(chan struct{})
	h.Shutdown = shutdown
	globalMutex.Unlock()
	go h.Reciever(h.ID, handlerNum)
	go h.Sender(h.ID, handlerNum, shutdown)
	log.Infof("Start a new stream! ID:%d  HandleNumber:%d direction:%d\n", h.ID, handlerNum, s.Stat().Direction)
}

//Reciever for NodeStreamHandler
func (h *NodeStreamHandler) Reciever(ID, handleNum int) {
	log.Infof("---Handler-%d-%d Reciever Start---", h.ID, handleNum)
	em, err := (*h.Host).EventBus().Emitter(new(NodeInfoMessage))
	if err != nil {
		panic(err)
	}
	defer em.Close()
	for {
		str, err := h.ReadWriter.ReadString('\n')
		if err != nil {
			log.Error(err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		log.Infof("#Reader%d read: %v\n", h.ID, str)
	}
}

//Sender for NodeStreamHandler
func (h *NodeStreamHandler) Sender(ID, handleNum int, shutdown <-chan struct{}) {
	log.Infof("---Handler-%d-%d Sender Start---", h.ID, handleNum)
	queue := Bus.Broadcast.Chan(-1)
	for {
		select {
		case <-shutdown:
			break
		case item := <-queue:
			log.Infof("#Sender%d Recieve TestQueue Message item:%v\n", h.ID, item)
		}
	}
}
