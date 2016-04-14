package main

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"sync"
	"sync/atomic"

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

func (s *Statistic) Get(name string) *Traffic {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.vars[name] == nil {
		s.vars[name] = &Traffic{Address: IPKey}
	}
	return s.vars[name]
}

func (t Traffics) Len() int {
	return len(t)
}

func (t Traffics) Less(i, j int) bool {
	return t[i].Inbound+t[i].Outbound < t[j].Inbound+t[j].Outbound
}

func (t Traffics) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (s Statistic) SetTraffic(dstIP, srcIP net.IP, dataLen uint64) {
	_isInbound := isInbound(localhost, dstIP)
	if _isInbound {
		IPKey = dstIP.String()
	} else {
		IPKey = srcIP.String()
	}

	traffic := s.Get(IPKey)
	if _isInbound {
		atomic.AddUint64(&traffic.Inbound, dataLen)
	} else {
		atomic.AddUint64(&traffic.Outbound, dataLen)
	}
}

func (s Statistic) PrintSortedStatisticString() {
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
	for _, v := range ts {
    sum = v.Inbound+v.Outbound
		fmt.Fprintf(
			&buf,
			"[%v] Traffic: %v / Inbound: %v / Outbound: %v\n",
			v.Address,
			bytefmt.ByteSize(sum),
			bytefmt.ByteSize(v.Inbound),
			bytefmt.ByteSize(v.Outbound),
		)
	}
	return buf.String()
}
