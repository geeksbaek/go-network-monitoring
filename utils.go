package main

import (
	"fmt"
	"net"
)

var (
	isInboundBool bool
	i             int
)

// ConsoleClear is a clear the console screen.
func ConsoleClear() {
	fmt.Print("\033[H\033[2J")
}

func lookupAddr(addr string) string {
	domain, err := net.LookupAddr(addr)
	if err != nil || len(domain) < 1 {
		return addr
	}
	return domain[0]
}

func isInbound(localhostPattern, dstIP []byte) bool {
	isInboundBool = false
	for i = range localhostPattern {
		if localhostPattern[i] == dstIP[i] {
			isInboundBool = true
		} else {
			return isInboundBool
		}
	}
	return isInboundBool
}
