package main

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"sync"
	"sync/atomic"
	"runtime"

	"github.com/pivotal-golang/bytefmt"
)

var (
	// localhost = []byte{172, 31}
	localhost = []byte{192, 168}
)

type Statistic struct {
	mutex *sync.RWMutex
	vars  map[uint32]*Traffic
}

type Traffic struct {
	Address  net.IP
	Inbound  uint64
	Outbound uint64
}

type Traffics []*Traffic

func (s *Statistic) Get(IP net.IP) *Traffic {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	IPKey := IPtoUint32(IP)
	if s.vars[IPKey] == nil {
		s.vars[IPKey] = &Traffic{Address: IP}
	}
	return s.vars[IPKey]
}

func (s *Statistic) SetTraffic(dstIP, srcIP net.IP, dataLen uint64) {
	if isInbound(localhost, dstIP) {
		atomic.AddUint64(&s.Get(dstIP).Inbound, dataLen)
	} else {
		atomic.AddUint64(&s.Get(srcIP).Outbound, dataLen)
	}
}

func (s *Statistic) PrintSortedStatisticString() {
	s.mutex.RLock()
	ts := make(Traffics, 0, len(s.vars))
	for _, t := range s.vars {
		ts = append(ts, t)
	}
	s.mutex.RUnlock()
	sort.Sort(sort.Reverse(ts))
	fmt.Print(ts.String())
}

func (ts Traffics) String() string {
	var buf bytes.Buffer
	var sum uint64
	
	buf.WriteString("\033[H\033[2J") // for clear the screen
	
	fmt.Fprintf(&buf, "Running Goroutines : %d\n", runtime.NumGoroutine())
	
	for _, v := range ts {
		sum = v.Inbound + v.Outbound
		fmt.Fprintf(
			&buf,
			"[%s] Traffic: %-7s / Inbound: %-7s / Outbound: %-7s\n",
			v.Address.String(),
			bytefmt.ByteSize(sum),
			bytefmt.ByteSize(v.Inbound),
			bytefmt.ByteSize(v.Outbound),
		)
	}
	return buf.String()
}
