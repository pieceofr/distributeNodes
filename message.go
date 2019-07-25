package main

import (
	"encoding/json"
)

type ExternalMessage interface {
	Marshal() ([]byte, error)
	Unmarshal(msgByte []byte) error
}

// NodeInfoMessage is a message to inform peers about noder
type NodeInfoMessage struct {
	NodeType  `json:"nodeType"`
	Address   string `json:"address"`
	PublicKey string `json:"publicKey"`
}

func (m *NodeInfoMessage) Marshal() ([]byte, error) {
	info, err := json.Marshal(m)
	return info, err
}

func (m *NodeInfoMessage) Unmarshal(msgByte []byte) error {
	err := json.Unmarshal(msgByte, m)
	return err
}
