package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bitmark-inc/logger"
)

var (
	cfg     config
	log     *logger.L
	theNode PeerNode
)

func main() {
	path := filepath.Join(os.Getenv("PWD"), "servant.conf")
	flag.StringVar(&path, "conf", "", "Specify configuration file")
	flag.Parse()
	if err := ParseConfigurationFile(path, &cfg); err != nil {
		panic(fmt.Sprintf("config file read failed: %s", err))
	}
	if err := logger.Initialise(cfg.Logging); err != nil {
		panic(fmt.Sprintf("logger initialization failed: %s", err))
	}
	log = logger.New("nodes")
	go theNode.MessageCenter()
	//Initialize messagebus
	messagebusInit()
	// Initialize a peer node
	if err := theNode.Init(cfg); err != nil {
		panic(fmt.Sprintf("node initialization failed: %s", err))
	}
	defer func() {
		log.Info("----- Reset and Exit -----")
		theNode.Reset()
	}()

	for {
		time.Sleep(10 * time.Second)
	}

}
