// TODO.
// stat map을 정렬하는 과정에서 key(ip address)의 도메인을 구한 뒤
// 해당 도메인 네임을 새로운 키로 사용하도록 한다.
// 복수의 ip가 동일한 도메인 네임을 사용하는 경우, 그것을 하나로 취급하기 위해서다.

package main

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/google/gopacket/layers"
	"github.com/pivotal-golang/bytefmt"
)

const (
	localhost = `192\.168.*`
)

type Statistic map[string]*Traffic

type Traffic struct {
	Address         string
	InboundTraffic  uint64
	OutboundTraffic uint64
	Traffic         uint64
}

type Traffics []*Traffic

func (t Traffics) Len() int {
	return len(t)
}

func (t Traffics) Less(i, j int) bool {
	return t[i].Traffic < t[j].Traffic
}

func (t Traffics) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (stats Statistic) SetTraffic(ip *layers.IPv4, t uint64) {
	dstIP := ip.DstIP.String()
	srcIP := ip.SrcIP.String()

	if isInbound(localhost, dstIP) {
		if stats[dstIP] == nil {
			stats[dstIP] = &Traffic{}
		}
		stats[dstIP].Address = dstIP
		stats[dstIP].InboundTraffic += t
		stats[dstIP].Traffic += t
	} else {
		if stats[srcIP] == nil {
			stats[srcIP] = &Traffic{}
		}
		stats[srcIP].Address = srcIP
		stats[srcIP].OutboundTraffic += t
		stats[srcIP].Traffic += t
	}
}

func (stats Statistic) SortedString() string {
	x := make(Traffics, 0, len(stats))
	for _, stat := range stats {
		x = append(x, stat)
	}
	sort.Sort(sort.Reverse(x))
	return x.String()
}

func (t Traffics) String() string {
	var buf bytes.Buffer
	buf.WriteString("\033[H\033[2J") // for clear the screen
	for _, v := range t {
		fmt.Fprintf(
			&buf,
			"[%v] Total Traffic: %v / Inbound: %v / Outbound: %v\n",
			lookupAddr(v.Address),
			bytefmt.ByteSize(v.Traffic),
			bytefmt.ByteSize(v.InboundTraffic),
			bytefmt.ByteSize(v.OutboundTraffic),
		)
	}
	return buf.String()
}
