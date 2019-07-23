package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
)

const (
	//MessageSeperator seperator between header and content
	MessageSeperator = "::"
	//HeaderAnnounceSelf s the header to announce self
	HeaderAnnounceSelf = "peer"
)

// NodeInfoMessage is a message to inform peers about noder
type NodeInfoMessage struct {
	NodeType
	Address   string
	PublicKey string
}

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
	NodeInfo   NodeInfoMessage
}

//Setup setup Handler
func (h *NodeStreamHandler) Setup(info NodeInfoMessage) {
	h.NodeInfo = info
}

// Handler  for streamHandler
func (h *NodeStreamHandler) Handler(s network.Stream) {
	log.Infof("Start a new stream! direction:%d\n", s.Stat().Direction)
	h.Stream = s
	h.ReadWriter = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go h.Reciever()
	go h.Sender()
}

//Reciever for NodeStreamHandler
func (h *NodeStreamHandler) Reciever() error {
	for {
		str, err := h.ReadWriter.ReadString('\n')
		if err != nil {
			log.Errorf("Read buff error:%s\n", err.Error())
			break
		}
		header, content, err := h.MessageParser(str)
		if err != nil {
			log.Errorf("MessageParser error : %s\n", err.Error())
			break
		}
		var peerInfo NodeInfoMessage
		switch header {
		case "peer":
			err := json.Unmarshal([]byte(content), &peerInfo)
			if err != nil {
				log.Error(errMessageFormat.Error())
				break
			}
			log.Debugf("peer message recieved:%s\n", peerInfo.Address)
			Bus.TestQueue.Send("peer", []byte(peerInfo.Address))
		default:
			log.Error(errMessageFormat.Error())
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

//Sender for NodeStreamHandler
func (h *NodeStreamHandler) Sender() {
	for {
		message, err := h.MessageComposer(HeaderAnnounceSelf, h.NodeInfo.String())
		if err != nil {
			break
		}
		if h.NodeInfo.NodeType != Client {
			h.ReadWriter.WriteString(fmt.Sprintf("%s\n", message))
			h.ReadWriter.Flush()
			log.Infof("Sender SEND: %s\n", fmt.Sprintf("%s\n", message))
		}
		time.Sleep(5 * time.Second)
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
