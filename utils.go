package main

import (
	"fmt"
	"net"
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

func isInbound(inboundPattern, dstIP []byte) bool {
	for i := range inboundPattern {
		if inboundPattern[i] != dstIP[i] {
			return false
		}
	}
	return true
}
