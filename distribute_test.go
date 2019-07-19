package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/bitmark-inc/logger"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	curPath := os.Getenv("PWD")
	var logLevel map[string]string
	logLevel = make(map[string]string, 0)
	logLevel["DEFAULT"] = "info"

	var logConfig = logger.Configuration{
		Directory: curPath,
		File:      "unittest.log",
		Size:      1048576,
		Count:     20,
		Console:   true,
		Levels:    logLevel,
	}
	if err := logger.Initialise(logConfig); err != nil {
		panic(fmt.Sprintf("logger initialization failed: %s", err))
	}
	log = logger.New("nodes")

	os.Exit(m.Run())
}

func TestIsSameNode(t *testing.T) {
	addr := `/ip4/118.163.120.180/tcp/12136/p2p/QmeHidiFxXLosH44wkAyq3modWku5gprV168t9VeH1JSyV`
	addrDiffPort := "/ip4/118.163.120.180/tcp/12140/p2p/QmeHidiFxXLosH44wkAyq3modWku5gprV168t9VeH1JSyV"
	addrDiffPublicPort := "/ip4/118.163.120.180/tcp/12140/p2p/QmeHidiFxXLosH44wkAyq3modWku5gprV168t9VeH1JSyV"
	addrDiffID := "/ip4/118.163.120.180/tcp/12136/p2p/QmeHidiFxXLosH44wkAyq3modWku5gprV168t9VeH1JSyV"
	var p = PeerNode{PublicIP: "118.163.120.180", Port: "12136"}
	assert.True(t, p.IsSameNode(addr), "Should be the same")
	assert.True(t, !p.IsSameNode(addrDiffPort), "Should not be the same")
	assert.True(t, !p.IsSameNode(addrDiffPublicPort), "Should not be the same")
	assert.True(t, p.IsSameNode(addrDiffID), "Should be the same")

}

func TestGetServer(t *testing.T) {
	for i := 0; i < 10; i++ {
		server, err := GetAServer()
		assert.NoError(t, err, "get server error")
		assert.NotEmpty(t, len(server), "Empty Server")
	}
}
func TestLoadServer(t *testing.T) {
	curPath := os.Getenv("PWD")
	servantPath := path.Join(curPath, "servant.addr")
	paths, err := LoadServer(servantPath)
	assert.NoError(t, err, "error to loading path")
	fmt.Println(paths)
}

func TestConfig(t *testing.T) {
	curPath := os.Getenv("PWD")
	confpath := path.Join(curPath, "servant.conf")
	fmt.Println("config file:", confpath)
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

func TestClientConnect(t *testing.T) {
	var servantCfg config
	var servantNode PeerNode
	path := filepath.Join(os.Getenv("PWD"), "servant.conf")
	if err := ParseConfigurationFile(path, &servantCfg); err != nil {
		panic(fmt.Sprintf("config file read failed: %s", err))
	}

	if err := servantNode.Init(servantCfg); err != nil {
		panic(fmt.Sprintf("node initialization failed: %s", err))
	}

	time.Sleep(3 * time.Second)
	var clientConfig config
	var clientNode PeerNode
	pathClient := filepath.Join(os.Getenv("PWD"), "client0.conf")
	if err := ParseConfigurationFile(pathClient, &clientConfig); err != nil {
		panic(fmt.Sprintf("config file read failed: %s", err))
	}

	if err := clientNode.Init(clientConfig); err != nil {
		panic(fmt.Sprintf("node initialization failed: %s", err))
	}

}
