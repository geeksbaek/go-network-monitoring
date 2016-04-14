package main

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"sync"

	"github.com/pivotal-golang/bytefmt"
)

var (
	IPKey     string
	localhost = []byte{192, 168}
)

type Statistic struct {
	mutex *sync.RWMutex
	vars  map[string]*Traffic
}

type Traffic struct {
	Address  string
	Inbound  uint64
	Outbound uint64
}

type Traffics []*Traffic

func (t Traffics) Len() int {
	return len(t)
}

func (t Traffics) Less(i, j int) bool {
	return t[i].Inbound+t[i].Outbound < t[j].Inbound+t[j].Outbound
}

func (t Traffics) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (stats Statistic) SetTraffic(dstIP, srcIP net.IP, traffic uint64) {
	stats.mutex.Lock()
	defer stats.mutex.Unlock()

	if isInbound(localhost, dstIP) {
		IPKey = dstIP.String()
		if stats.vars[IPKey] == nil {
			stats.vars[IPKey] = &Traffic{}
		}
		stats.vars[IPKey].Inbound += traffic
	} else {
		IPKey = srcIP.String()
		if stats.vars[IPKey] == nil {
			stats.vars[IPKey] = &Traffic{}
		}
		stats.vars[IPKey].Outbound += traffic
	}
	stats.vars[IPKey].Address = IPKey
}

func (stats Statistic) PrintSortedStatisticString() {
	stats.mutex.Lock()
	x := make(Traffics, 0, len(stats.vars))
	for _, stat := range stats.vars {
		x = append(x, stat)
	}
	defer stats.mutex.Unlock()
	sort.Sort(sort.Reverse(x))
	fmt.Print(x.String())
}

func (t Traffics) String() string {
	var buf bytes.Buffer
	buf.WriteString("\033[H\033[2J") // for clear the screen

	for _, v := range t {
		fmt.Fprintf(
			&buf,
			"[%v] Total Traffic: %v / Inbound: %v / Outbound: %v\n",
			lookupAddr(v.Address),
			bytefmt.ByteSize(v.Inbound+v.Outbound),
			bytefmt.ByteSize(v.Inbound),
			bytefmt.ByteSize(v.Outbound),
		)
	}
	return buf.String()
}
