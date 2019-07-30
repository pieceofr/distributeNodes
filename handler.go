package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
)

const (
	//MessageSeperator seperator between header and content
	MessageSeperator = "::"
	//HeaderAnnounceSelf s the header to announce self
	HeaderAnnounceSelf = "peer"
	cycleInterval      = 10 * time.Second
)

//BiStreamHandler Bidirectional  StreamHandler
type BiStreamHandler interface {
	Handler(s network.Stream)
	Reciever() error
	Sender() error
}

//NodeStreamHandler for  node to handle stream
type NodeStreamHandler struct {
	Stream     network.Stream
	ReadWriter *bufio.ReadWriter
	NodeInfoMessage
	Shutdown chan<- struct{}
}

//Setup setup Handler
func (h *NodeStreamHandler) Setup(info NodeInfoMessage) {
	h.NodeInfoMessage = info
}

// Handler  for streamHandler
func (h *NodeStreamHandler) Handler(s network.Stream) {
	log.Infof("Start a new stream! direction:%d\n", s.Stat().Direction)
	h.Stream = s
	h.ReadWriter = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	shutdown := make(chan struct{})
	h.Shutdown = shutdown
	go h.Reciever()
	go h.Sender(shutdown)
}

//Reciever for NodeStreamHandler
func (h *NodeStreamHandler) Reciever() error {
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
			log.Infof("RECIEVE From %v Extra:%v", shortID(peerInfo.ID), peerInfo.Extra)
			Bus.TestQueue.Send("peer", []byte(fmt.Sprintf("%v", peerInfo.NodeType)), []byte(peerInfo.ID), []byte(peerInfo.Address))
		default:
			log.Error(errMessageFormat.Error())
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

//Sender for NodeStreamHandler
func (h *NodeStreamHandler) Sender(shutdown <-chan struct{}) {
	queue := Bus.Broadcast.Chan(-1)
	cycleTimer := time.After(cycleInterval)
	for {
		select {
		case <-shutdown:
			break
		case item := <-queue:
			switch item.Command {
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
				peerInfo.Extra = time.Now().UTC().String()
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
				log.Infof("Broadcasting Item Send: %s \n", message)
				h.ReadWriter.WriteString(fmt.Sprintf("%s\n", message))
				h.ReadWriter.Flush()
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
				log.Debugf("Broadcasting Self ID:%s type:%d\n", shortID(h.NodeInfoMessage.ID), h.NodeInfoMessage.NodeType)
			}
		}
		//log.Infof("NodeType:%v Sender  %s SEND", h.NodeInfoMessage, h.NodeInfoMessage.PublicKey[len(h.NodeInfoMessage.PublicKey)-10:len(h.NodeInfoMessage.PublicKey)-1])
		time.Sleep(10 * time.Second)
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
