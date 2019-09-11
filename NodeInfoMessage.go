package main

import (
	"encoding/json"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
)

// NodeInfoMessage is a message to inform peers about noder
type NodeInfoMessage struct {
	NodeType `json:"nodeType"`
	Address  []string `json:"address"`
	ID       string   `json:"id"`
	Extra    string
}

//Marshal  convert NodeInfoMessage to []byte
func (m *NodeInfoMessage) Marshal() ([]byte, error) {
	info, err := json.Marshal(m)
	return info, err
}

//Unmarshal  convert byte[] to NodeInfoMessage
func (m *NodeInfoMessage) Unmarshal(msgByte []byte) error {
	err := json.Unmarshal(msgByte, m)
	return err
}

// addrToConnAddr remove protocol ID and node ID
func addrToConnAddr(addr string) string {
	addrSlice := strings.Split(addr, "/")
	var retAddr string
	if len(addrSlice) > 4 {
		for idx, addr := range addrSlice[:len(addrSlice)-2] {
			if idx == (len(addrSlice) - 3) {
				retAddr = retAddr + addr
			} else {
				retAddr = retAddr + addr + "/"
			}
		}
	}
	return retAddr
}

func shortID(id string) string {
	ID, err := peer.IDFromString(id)
	if err != nil {
		return ""
	}
	return ID.ShortString()
}
