package main

import (
	"bufio"
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
	Identity
}

//Setup setup Handler
func (h *NodeStreamHandler) Setup(ID Identity) {
	// Create a buffer stream for non blocking read and write.
	h.Identity = ID
}

// Handler  for streamHandler
func (h *NodeStreamHandler) Handler(s network.Stream) {
	log.Infof("Start a new stream! direction:%d\n", s.Stat().Direction)
	h.Stream = s
	h.ReadWriter = bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go h.Reciever()
	go h.Sender()

	// stream 's' will stay open until you close it (or the other side closes it).
}

//Reciever for NodeStreamHandler
func (h *NodeStreamHandler) Reciever() error {
	for {
		str, _ := h.ReadWriter.ReadString('\n')
		header, content, err := h.MessageParser(str)
		if err != nil {
			log.Error(err.Error())
		}
		switch header {
		case "peer":
			log.Infof("header:%s\ncontent:%s", header, content)
		default:
			log.Error(errMessageFormat.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

//Sender for NodeStreamHandler
func (h *NodeStreamHandler) Sender() {
	for {
		ID, err := h.Identity.MarshalPrvKey()
		if err != nil {
			log.Error(err.Error())
			break
		}
		message, err := h.MessageComposer(HeaderAnnounceSelf, ID)
		if err != nil {
			break
		}
		h.ReadWriter.WriteString(fmt.Sprintf("%s\n", message))
		h.ReadWriter.Flush()
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
func (h *NodeStreamHandler) MessageComposer(header string, content string) (message string, err error) {
	return fmt.Sprintf("%s::%s", header, content), nil
}
