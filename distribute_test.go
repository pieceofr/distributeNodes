package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	curPath := os.Getenv("PWD")
	confpath := path.Join(curPath, "servant.conf")
	fmt.Println("confile:", confpath)
	var cfg config
	err := ParseConfigurationFile(confpath, &cfg)
	assert.NoError(t, err, "ParseConfig Error")
	fmt.Println("NodeType:", cfg.NodeType)
	fmt.Println("Public IP:", cfg.PublicIP)
	fmt.Println("Port:", cfg.Port)
	fmt.Println("use static:", cfg.StaticIdentity.UseStatic)
	fmt.Println("private key:", cfg.StaticIdentity)
}

func TestMarshalUnmarshal(t *testing.T) {
	prv, err := randKey()
	//Marshall Key
	assert.NoError(t, err, "randKey Error")
	id := Identity{PrvKey: prv}
	serialKey, err := id.MarshalPrvKey()
	assert.NoError(t, err, "Marshal Private Key Error")

	err = id.UnmarshalPrvKey(serialKey)
	assert.NoError(t, err, "Unmarshal Private Key Error")
	serialKey2, err := id.MarshalPrvKey()
	assert.NoError(t, err, "Marshal Private Key2 Error")
	assert.Equal(t, serialKey, serialKey2, "Error on Marshal and Unmarshal keys")
}

func TestGenIdentity(t *testing.T) {
	for i := 0; i < 1; i++ {
		var id Identity
		err := id.randIdentity()
		if err != nil {
			fmt.Println("randIdentity err:", err)
		}
		assert.NoError(t, err, "Generate identity error i="+strconv.Itoa(i))
		assert.NotNil(t, id.PrvKey)
		keyLen, err := id.PrvKey.Bytes()
		assert.NoError(t, err, "get key byte error")
		assert.NotZero(t, len(keyLen), " is zero")
	}
}

func TestSaveLoadPrivateKey(t *testing.T) {
	var node PeerNode
	err := node.NewRandomNode()
	assert.NoError(t, err, "gen random node error")
	assert.NotNil(t, node.Identity.PrvKey, "generate Identity fail")
	oriKey, err := node.MarshalPrvKey()
	assert.NoError(t, err, "marshal key error")
	node.SaveIdentity()
	var newNode PeerNode
	newNode.LoadIdentity()
	newKey, err := newNode.MarshalPrvKey()
	assert.NoError(t, err, "marshal key error")
	assert.Equal(t, oriKey, newKey, "oriKey, newKey not equal")

}
