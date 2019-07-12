package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

//Identity A identity who own this node
type Identity struct {
	PrvKey crypto.PrivKey
	Port   int
}

func randKey() (crypto.PrivKey, error) {
	r := rand.Reader
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	if err != nil {
		return nil, err
	}
	return prvKey, nil
}

//PublicKey get the public key of identity
func (i *Identity) PublicKey() (crypto.PubKey, error) {
	if nil == i.PrvKey {
		return nil, errors.New("private key is not initialized")
	}
	publicKey := i.PrvKey.GetPublic()
	if nil == publicKey {
		return nil, errors.New("generate public key error")
	}
	return publicKey, nil
}

func (i *Identity) randIdentity() error {
	prv, err := randKey()
	if err != nil {
		return err
	}
	i.PrvKey = prv
	//	fmt.Println("private key:", prvByte, " public key:", pubByte)
	// Bitmark Open Port from 12130-12150
	// 12130-12136 reserve for servant node
	// random port open from 12137 - 12150
	b := make([]byte, 1)
	_, err = rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	i.Port = 12137 + int(math.Mod(float64(b[0]), float64(13)))
	return nil
}
