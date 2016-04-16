package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// vars for decode
var (
	eth       layers.Ethernet
	ip4       layers.IPv4
	parser    *gopacket.DecodingLayerParser
	decoded   = []gopacket.LayerType{}
	packet    gopacket.Packet
	data      []byte
	layerType gopacket.LayerType
)

// vars for statistic
var (
	ticker    = time.Tick(time.Second * 2)
	statistic = &Statistic{
		mutex: new(sync.RWMutex),
		vars:  make(map[uint32]*Traffic),
	}
)

func Sniff(packetChannel <-chan gopacket.Packet) {
	parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4)
	for {
		select {
		case packet = <-packetChannel:
			gotPacket()
		case <-ticker:
			go fmt.Print(statistic.SortedStatisticString())
		}
	}
}

func gotPacket() {
	data = packet.Data()
	parser.DecodeLayers(data, &decoded)
	for _, layerType = range decoded {
		switch layerType {
		case layers.LayerTypeIPv4:
			go statistic.SetTraffic(ip4.DstIP, ip4.SrcIP, uint64(len(data)))
		default:
			atomic.AddUint64(&statistic.total, uint64(len(data)))
		}
	}
}
