package main

import (
	"errors"
	"fmt"
	"time"
)

type messageRecord struct {
	messages map[string][]string
}

func (m messageRecord) init() error {
	m.messages = make(map[string][]string)
	return nil
}

func (m messageRecord) addMessage(msg string) error {
	recievedTime := fmt.Sprintf("%d", time.Now().Unix())
	i, ok := m.messages[msg]
	if ok {
		m.messages[msg] = append(i, recievedTime)
	}
	return nil
}

func (m messageRecord) Messages() (map[string][]string, error) {
	if m.messages != nil {
		return m.messages, nil
	}
	return nil, errors.New("Not Inititialized")
}
