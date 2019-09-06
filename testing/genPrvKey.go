package main

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	crypto "github.com/libp2p/go-libp2p-crypto"
)

func main() {
	//Chnage here to generage key
	//count ->postfix and max number
	//prefix -> name of file
	// go run genPrvKey.go
	for i := 1; i <= 20; i++ {
		saveGenKey(i, "servant")
	}

}

func saveGenKey(count int, prefix string) error {
	prv, err := randKey()
	if err != nil {
		return err
	}

	keyfile := path.Join(os.Getenv("PWD"), "key", prefix+strconv.Itoa(count)+".prv")

	encodedKey, err := marshalPrvKey(prv)
	if err := ioutil.WriteFile(keyfile, []byte(encodedKey), 0644); err != nil {
		return err
	}

	return nil
}

func randKey() (crypto.PrivKey, error) {
	r := rand.Reader
	prvKey, _, err := crypto.GenerateEd25519Key(r)
	if err != nil {
		return nil, err
	}
	return prvKey, nil
}
func marshalPrvKey(prvKey crypto.PrivKey) ([]byte, error) {
	marshalKey, err := crypto.MarshalPrivateKey(prvKey)
	if err != nil {
		return nil, err
	}
	hexEncodeKey := make([]byte, hex.EncodedLen(len(marshalKey)))
	hex.Encode(hexEncodeKey, marshalKey)
	return hexEncodeKey, nil
}
