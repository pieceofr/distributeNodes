package main

import (
	"errors"
	"fmt"
)

var (
	errSameNode           = errors.New("The Same Peer Node")
	errMessageFormat      = errors.New("Message Format Is Invalid")
	errPeerHasInPeerStore = errors.New("peer has in the peer store")
	errReconnectStream    = errors.New("Stream Recconect Fail")
)

//ErrCombind all Errors together
func ErrCombind(errs ...error) (finalErr error) {
	for _, err := range errs {
		if err != nil {
			if finalErr != nil {
				finalErr = fmt.Errorf("%s-%s", finalErr, err)
			} else {
				finalErr = fmt.Errorf("%s", err)
			}
		}
	}
	return finalErr
}
