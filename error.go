package main

import (
	"errors"
	"fmt"
)

var (
	//ErrSameNode the same peer node
	ErrSameNode = errors.New("The Same Peer Node")
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
