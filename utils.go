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
	if len(dstIP) != 4 {
		return false
	}
	
	for i := range inboundPattern {
		if inboundPattern[i] != dstIP[i] {
			return false
		}
	}
	return true
}

func IPtoUint32(ip []byte) uint32 {
	if len(ip) != 4 {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// For Traffics sorting
func (t Traffics) Len() int {
	return len(t)
}

func (t Traffics) Less(i, j int) bool {
	return t[i].Inbound+t[i].Outbound < t[j].Inbound+t[j].Outbound
}

func (t Traffics) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
