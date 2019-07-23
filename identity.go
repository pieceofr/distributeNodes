package main

import (
	"crypto/rand"
	"errors"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mr-tron/base58"
)

//Identity A identity who own this node
type Identity struct {
	PrvKey crypto.PrivKey
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
	return nil
}

//MarshalPrvKey from base58 string to private key
func (i *Identity) MarshalPrvKey() (string, error) {
	marshalKey, err := crypto.MarshalPrivateKey(i.PrvKey)
	if err != nil {
		return "", err
	}
	encoded := base58.Encode(marshalKey)
	return encoded, nil
}

//UnmarshalPrvKey from base58 string to private key
func (i *Identity) UnmarshalPrvKey(prvKey string) error {
	decoded, err := base58.Decode(prvKey)
	if err != nil {
		return err
	}
	unmarshalKey, err := crypto.UnmarshalPrivateKey(decoded)
	if err != nil {
		return err
	}
	i.PrvKey = unmarshalKey
	return nil
}

//MarshalPublicKey from base58 string to public  key
func (i *Identity) MarshalPublicKey() (string, error) {
	pkey, err := i.PublicKey()
	if err != nil {
		return "", err
	}
	marshalKey, err := crypto.MarshalPublicKey(pkey)
	if err != nil {
		return "", err
	}
	encoded := base58.Encode(marshalKey)
	return encoded, nil
}

//UnmarshalPublicKey from base58 string to private key
func (i *Identity) UnmarshalPublicKey(pubKey string) (crypto.PubKey, error) {
	decoded, err := base58.Decode(pubKey)
	if err != nil {
		return nil, err
	}
	unmarshalKey, err := crypto.UnmarshalPublicKey(decoded)
	if err != nil {
		return nil, err
	}

	return unmarshalKey, nil
}
