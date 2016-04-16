package main

import (
	"fmt"
	"net"
)

const (
	MIDDLE_LEFT = rune('├')
	BOTTOM_LEFT = rune('└')
)

func GetBoxDrawingChar(i, maxDisplayLine, lenSlice int) rune {
	if i == maxDisplayLine || i == lenSlice-1 {
		return BOTTOM_LEFT
	}
	return MIDDLE_LEFT
}

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

func MeAndYou(dstIP, srcIP []byte) (net.IP, net.IP) {
	for i := range localhost {
		if localhost[i] != dstIP[i] {
			return srcIP, dstIP
		}
	}
	return dstIP, srcIP
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

func (t ConnTraffics) Len() int {
	return len(t)
}

func (t ConnTraffics) Less(i, j int) bool {
	return t[i].Inbound+t[i].Outbound < t[j].Inbound+t[j].Outbound
}

func (t ConnTraffics) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
