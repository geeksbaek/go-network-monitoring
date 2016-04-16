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

const (
	maxDisplayLine = 5
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
	Domain   string
	Inbound  uint64
	Outbound uint64
}

type (
	Traffics     []*Traffic
	ConnTraffics []*ConnTraffic
)

func (s *Statistic) Get(IP net.IP) *Traffic {
	IPKey := IPtoUint32(IP)
	
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.vars[IPKey] == nil {
		s.vars[IPKey] = &Traffic{
			Address: IP,
			MostConn: &ConnStatistic{
				mutex: new(sync.RWMutex),
				vars:  make(map[uint32]*ConnTraffic, 100),
			},
		}
	}
	return s.vars[IPKey]
}

func (s *ConnStatistic) Get(IP net.IP) *ConnTraffic {
	IPKey := IPtoUint32(IP)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.vars[IPKey] == nil {
		s.vars[IPKey] = &ConnTraffic{Address: IP}
	}
	return s.vars[IPKey]
}

func (s *Statistic) SetTraffic(dstIP, srcIP net.IP, dataLen uint64) {
	if me, you := MeAndYou(dstIP, srcIP); me.Equal(dstIP) {
		atomic.AddUint64(&s.Get(me).Inbound, dataLen)
		atomic.AddUint64(&s.Get(me).MostConn.Get(you).Outbound, dataLen)
	} else {
		atomic.AddUint64(&s.Get(me).Outbound, dataLen)
		atomic.AddUint64(&s.Get(me).MostConn.Get(you).Inbound, dataLen)
	}
}

func (s *Statistic) SortedStatisticString() string {
	s.mutex.RLock()
	ts := make(Traffics, 0, len(s.vars))
	for _, t := range s.vars {
		ts = append(ts, t)
	}
	s.mutex.RUnlock()
	sort.Sort(sort.Reverse(ts))
	return ts.String()
}

func (ts Traffics) String() string {
	var buf bytes.Buffer

	buf.WriteString("\033[H\033[2J") // for clear the screen
	fmt.Fprintf(&buf, "Running Goroutines : %d / Total Traffic : %s\n\n",
		runtime.NumGoroutine(),
		bytefmt.ByteSize(atomic.LoadUint64(&statistic.total)),
	)

	for _, v := range ts {
		fmt.Fprintf(
			&buf,
			"[%s] Traffic: %s / Inbound: %s / Outbound: %s\n%s\n",
			v.Address.String(),
			bytefmt.ByteSize(v.Inbound+v.Outbound),
			bytefmt.ByteSize(v.Inbound),
			bytefmt.ByteSize(v.Outbound),
			v.MostConn.SortedStatisticString(),
		)
	}
	return buf.String()
}

func (s *ConnStatistic) SortedStatisticString() string {
	s.mutex.RLock()
	ts := make(ConnTraffics, 0, len(s.vars))
	for _, t := range s.vars {
		ts = append(ts, t)
	}
	s.mutex.RUnlock()
	sort.Sort(sort.Reverse(ts))
	return ts.String()
}

func (ts ConnTraffics) String() string {
	var buf bytes.Buffer

	for i, v := range ts {
		if i > maxDisplayLine {
			break
		}

		fmt.Fprintf(
			&buf,
			"%c[%s] Traffic: %s / Inbound: %s / Outbound: %s\n",
			GetBoxDrawingChar(i, maxDisplayLine, len(ts)),
			//v.Domain,
			lookupAddr(v.Address.String()),
			bytefmt.ByteSize(v.Inbound+v.Outbound),
			bytefmt.ByteSize(v.Inbound),
			bytefmt.ByteSize(v.Outbound),
		)
	}
	return buf.String()
}
