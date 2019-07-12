package main

import (
	"errors"
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenIdentity(t *testing.T) {
	for i := 0; i < 100; i++ {
		var id Identity
		err := id.randIdentity()
		if err != nil {
			log.Println("randIdentity err:", err)
		}
		/*
			keyB, err := id.PrvKey.Bytes()
			fmt.Println("prvKey Gen:", keyB, " err:", err)
		*/
		assert.NoError(t, err, "Generate identity error i="+strconv.Itoa(i))
		assert.NotNil(t, id.PrvKey)
		keyLen, err := id.PrvKey.Bytes()
		assert.NoError(t, err, "get key byte error")
		assert.NotZero(t, len(keyLen), " is zero")
		if id.Port > 12150 || id.Port < 12137 {
			err = errors.New("port number is not in the range")
			assert.NoError(t, err, err)
		}
	}
}
