package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

// PeerWriters write to current Stream
var PeerWriters []*bufio.ReadWriter
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
	PeerWriters = append(PeerWriters, h.ReadWriter)
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
		header, content, err := h.MessageParser(str)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		switch header {
		case "peer":
			var peerInfo NodeInfoMessage
			err := peerInfo.Unmarshal([]byte(content))
			if err != nil {
				log.Error(ErrCombind(err, errMessageFormat).Error())
				time.Sleep(1 * time.Second)
				continue
			}
			log.Infof("READ -- Handler-%d-%d PeerInfo:%v Extra:%v", h.ID, handleNum, shortID(peerInfo.ID), peerInfo.Extra)
			//Bus.TestQueue.Send("peer", []byte(fmt.Sprintf("%v", peerInfo.NodeType)), []byte(peerInfo.ID), []byte(peerInfo.Address), []byte(peerInfo.Extra))
			em.Emit(peerInfo)
			if err != nil {
				log.Error(err.Error())
			}
		default:
			log.Error(errMessageFormat.Error())
		}
	}
}

//Sender for NodeStreamHandler
func (h *NodeStreamHandler) Sender(ID, handleNum int, shutdown <-chan struct{}) {
	log.Infof("---Handler-%d-%d Sender Start---", h.ID, handleNum)
	queue := Bus.Broadcast.Chan(-1)
	cycleTimer := time.After(cycleInterval)
	for {
		select {
		case <-shutdown:
			break
		case item := <-queue:
			switch item.Command {
			case "testing":
				log.Infof("#Handler: %d Recieve Broadcasting Message\n", h.ID)
			case "peer":
				var peerInfo NodeInfoMessage
				nType, err := strconv.Atoi(string(item.Parameters[0]))
				if err != nil {
					log.Error(err.Error())
					continue
				}
				peerInfo.NodeType = NodeType(nType)
				peerInfo.ID = string(item.Parameters[1])
				peerInfo.Address = string(item.Parameters[2])
				//peerInfo.Extra = fmt.Sprintf("%v", time.Now())
				peerInfo.Extra = string(item.Parameters[3])
				infoOut, err := peerInfo.Marshal()
				if err != nil {
					log.Error(err.Error())
					continue
				}
				message, err := h.MessageComposer(HeaderAnnounceSelf, string(infoOut))
				if err != nil {
					time.Sleep(10 * time.Second)
					continue
				}
				for _, peerWriter := range PeerWriters {
					peerWriter.WriteString(fmt.Sprintf("%s\n", message))
					peerWriter.Flush()
					log.Infof("##PeerHandler-%d-%d Send Peers Announce Peer Info:%s :%v\n", h.ID, handleNum, peerInfo.Address)
				}
			}
		case <-cycleTimer:
			cycleTimer = time.After(cycleInterval)
			if h.NodeInfoMessage.NodeType != Client {
				h.NodeInfoMessage.Extra = time.Now().UTC().String()
				info, err := h.NodeInfoMessage.Marshal()
				if err != nil {
					log.Errorf("Marshal Error:%s\n", err)
					time.Sleep(5 * time.Second)
					continue
				}
				message, err := h.MessageComposer(HeaderAnnounceSelf, string(info))
				if err != nil {
					log.Errorf("MessageComposer Error:%s\n", err)
					time.Sleep(5 * time.Second)
					continue
				}

				h.ReadWriter.WriteString(fmt.Sprintf("%s\n", message))
				h.ReadWriter.Flush()
				log.Infof("Handler -%d-%d: Broadcasting Self ID:%s type:%d\n", h.ID, handlerNum, shortID(h.NodeInfoMessage.ID), h.NodeInfoMessage.NodeType)

				for _, peerWriter := range PeerWriters {
					peerWriter.WriteString(fmt.Sprintf("%s\n", message))
					peerWriter.Flush()
					log.Infof("PeerHandler-%d-%d Broadcasting Self ID:%s type:%d\n", h.ID, handleNum, shortID(h.NodeInfoMessage.ID), h.NodeInfoMessage.NodeType)
				}

			}
		}
		//log.Infof("NodeType:%v Sender  %s SEND", h.NodeInfoMessage, h.NodeInfoMessage.PublicKey[len(h.NodeInfoMessage.PublicKey)-10:len(h.NodeInfoMessage.PublicKey)-1])
		//time.Sleep(10 * time.Second)
	}
}

//MessageParser parsing messages
func (h *NodeStreamHandler) MessageParser(msg string) (header string, content string, err error) {
	s := strings.SplitN(msg, MessageSeperator, 2)
	if len(s) == 0 {
		return "", "", errMessageFormat
	} else if len(s) == 1 {
		return s[0], "", nil
	}
	return s[0], s[1], nil
}

//MessageComposer parsing messages
func (h *NodeStreamHandler) MessageComposer(header string, messages ...string) (message string, err error) {
	var content string
	for _, msg := range messages {
		if len(msg) != 0 {
			content = content + msg
		}
	}
	return fmt.Sprintf("%s::%s", header, content), nil
}

func (i *NodeInfoMessage) String() string {
	byteStr, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(byteStr)
}

func (h *NodeStreamHandler) emmitEvent(event interface{}) {
	switch event.(type) {
	case NodeInfoMessage:
		bus := (*h.Host).EventBus()
		em, err := bus.Emitter(event.(NodeInfoMessage))
		if err != nil {
			log.Errorf("emit event error:%v\n", err)
		}
		defer em.Close()
	}
}
