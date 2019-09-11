package main

import (
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
)

// addrToConnAddr remove protocol ID and node ID
func addrToConnAddr(addr string) string {
	addrSlice := strings.Split(addr, "/")
	var retAddr string
	if len(addrSlice) > 4 {
		for idx, addr := range addrSlice[:len(addrSlice)-2] {
			if idx == (len(addrSlice) - 3) {
				retAddr = retAddr + addr
			} else {
				retAddr = retAddr + addr + "/"
			}
		}
	}
	return retAddr
}

func shortID(id string) string {
	ID, err := peer.IDFromString(id)
	if err != nil {
		return ""
	}
	return ID.ShortString()
}
