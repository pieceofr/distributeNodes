package main

import (
	"encoding/json"
)

// NodeInfoMessage is a message to inform peers about noder
type NodeInfoMessage struct {
	NodeType `json:"nodeType"`
	Address  string `json:"address"`
	ID       string `json:"id"`
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

func shortID(id string) string {
	if len(id) > 11 {
		return id[len(id)-11 : len(id)-1]
	}
	return id
}
