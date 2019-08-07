package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bitmark-inc/logger"
)

var (
	cfg         config
	log         *logger.L
	theNode     PeerNode
	globalMutex *sync.Mutex
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
	globalMutex = &sync.Mutex{}
	//Initialize messagebus
	messagebusInit()
	// Initialize a peer node
	if err := theNode.setup(cfg); err != nil {
		panic(fmt.Sprintf("node initialization failed: %s", err))
	}
	go theNode.run()

	defer func() {
		log.Info("----- Reset and Exit -----")
		theNode.Reset()
	}()

	for {
		time.Sleep(10 * time.Second)
	}

}
