package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

//Identity A identity who own this node
type Identity struct {
	PrvKey crypto.PrivKey
}

func randKey() (crypto.PrivKey, error) {
	r := rand.Reader
	prvKey, _, err := crypto.GenerateEd25519Key(r)
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
	return nil
}

//MarshalPrvKey from hex encoded string to private key
func (i *Identity) MarshalPrvKey() ([]byte, error) {
	marshalKey, err := crypto.MarshalPrivateKey(i.PrvKey)
	if err != nil {
		return nil, err
	}
	hexEncodeKey := make([]byte, hex.EncodedLen(len(marshalKey)))
	hex.Encode(hexEncodeKey, marshalKey)
	return hexEncodeKey, nil
}

//UnmarshalPrvKey from hex string to private key
func (i *Identity) UnmarshalPrvKey(prvKey []byte) error {
	hexDecodeKey := make([]byte, hex.DecodedLen(len(prvKey)))
	_, err := hex.Decode(hexDecodeKey, prvKey)
	if err != nil {
		return err
	}
	unmarshalKey, err := crypto.UnmarshalPrivateKey(hexDecodeKey)
	if err != nil {
		return err
	}
	i.PrvKey = unmarshalKey
	return nil
}
