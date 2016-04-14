package main

import (
	"time"

	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// vars for decode
var (
	eth       layers.Ethernet
	ip4       layers.IPv4
	ip6       layers.IPv6
	parser    *gopacket.DecodingLayerParser
	decoded   = []gopacket.LayerType{}
	packet    gopacket.Packet
	data      []byte
	layerType gopacket.LayerType
)

// vars for statistic
var (
	ticker    = time.Tick(time.Second * 1)
	statistic = make(Statistic, 100000)
	rwMutex   = new(sync.RWMutex)
	traffic   uint64
)

func Sniff(packetChannel <-chan gopacket.Packet) {
	parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6)
	for {
		select {
		case packet = <-packetChannel:
			data = packet.Data()
			parser.DecodeLayers(data, &decoded)
			for _, layerType = range decoded {
				switch layerType {
				case layers.LayerTypeIPv6:
					traffic = uint64(len(data))
					go statistic.SetTraffic(ip6.DstIP, ip6.SrcIP, traffic)
				case layers.LayerTypeIPv4:
					traffic = uint64(len(data))
					go statistic.SetTraffic(ip4.DstIP, ip4.SrcIP, traffic)
				}
			}
		case <-ticker:
			go statistic.PrintSortedString()
		}
	}
}
