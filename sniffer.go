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
	parser    *gopacket.DecodingLayerParser
	decoded   = []gopacket.LayerType{}
	packet    gopacket.Packet
	data      []byte
	layerType gopacket.LayerType
)

// vars for statistic
var (
	ticker    = time.Tick(time.Second * 2)
	statistic = &Statistic{new(sync.RWMutex), make(map[string]*Traffic, 10000000)}
	dataLen   uint64
)

func Sniff(packetChannel <-chan gopacket.Packet) {
	parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4)
	for {
		select {
		case packet = <-packetChannel:
			data = packet.Data()
			parser.DecodeLayers(data, &decoded)
			for _, layerType = range decoded {
				switch layerType {
				case layers.LayerTypeIPv4:
					dataLen = uint64(len(data))
					go statistic.SetTraffic(ip4.DstIP, ip4.SrcIP, dataLen)
				}
			}
		case <-ticker:
			go statistic.PrintSortedStatisticString()
		}
	}
}
