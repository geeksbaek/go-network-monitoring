package main

import (
	_ "fmt"
	"net"
	"regexp"
)

// ConsoleClear is a clear the console screen.
// func ConsoleClear() {
//   fmt.Print("\033[H\033[2J")
// }

func lookupAddr(addr string) string {
	domain, err := net.LookupAddr(addr)
	if err != nil || len(domain) < 1 {
		return addr
	}
	return domain[0]
}

func isInbound(localhost, dstIP string) bool {
	if matched, _ := regexp.MatchString(localhost, dstIP); matched {
		return true
	}
	return false
}
