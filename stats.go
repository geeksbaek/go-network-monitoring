package main

import (
	"bytes"
	"fmt"
	"net"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/pivotal-golang/bytefmt"
)

var (
	// localhost = []byte{172, 31}
	localhost = []byte{192, 168}
)

type Statistic struct {
	mutex *sync.RWMutex
	vars  map[uint32]*Traffic
	total uint64
}

type Traffic struct {
	Address  net.IP
	Inbound  uint64
	Outbound uint64
	MostConn *ConnStatistic
}

type ConnStatistic struct {
	mutex *sync.RWMutex
	vars  map[uint32]*ConnTraffic
}

type ConnTraffic struct {
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
		s.vars[IPKey] = &Traffic{
			Address: IP,
			MostConn: &ConnStatistic{
				mutex: new(sync.RWMutex),
				vars:  make(map[uint32]*ConnTraffic),
			},
		}
	}
	return s.vars[IPKey]
}

func (s *ConnStatistic) Get(IP net.IP) *ConnTraffic {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	IPKey := IPtoUint32(IP)
	if s.vars[IPKey] == nil {
		s.vars[IPKey] = &ConnTraffic{
			Address: IP,
		}
	}
	return s.vars[IPKey]
}

func (s *Statistic) SetTraffic(dstIP, srcIP net.IP, dataLen uint64) {
	me, you := MeAndYou(dstIP, srcIP)

	if me.Equal(dstIP) {
		atomic.AddUint64(&s.Get(me).Inbound, dataLen)
		atomic.AddUint64(&s.Get(me).MostConn.Get(you).Outbound, dataLen)
	} else {
		atomic.AddUint64(&s.Get(me).Outbound, dataLen)
		atomic.AddUint64(&s.Get(me).MostConn.Get(you).Inbound, dataLen)
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
