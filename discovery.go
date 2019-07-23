package main

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"path"
	"time"
)

//LoadServer is to load servant from a file
func LoadServer(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer file.Close()
	var address []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) != 0 {
			address = append(address, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return address, nil
}

//GetAServer return a servant/server to connect
func GetAServer() (string, error) {
	curPath := os.Getenv("PWD")
	servantPath := path.Join(curPath, "servant.addr")
	servers, err := LoadServer(servantPath)
	if err != nil {
		return "", err
	}
	if len(servers) == 0 {
		return "", errors.New("no server availible")
	}
	sNum := rand.NewSource(time.Now().UnixNano())
	rNum := rand.New(sNum)

	if len(servers) > 1 {
		serverIdx := rNum.Intn(len(servers))
		return servers[serverIdx], nil
	}
	return servers[0], nil
}
