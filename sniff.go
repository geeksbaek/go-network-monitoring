package main

import (
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// vars for decode
var (
	eth     layers.Ethernet
	ip4     layers.IPv4
	ip6     layers.IPv6
	parser  *gopacket.DecodingLayerParser
	decoded = []gopacket.LayerType{}
)

// vars for statistic
var flow = &Flow{
	Vars:  map[Endpoints]Traffic{},
	Mutex: new(sync.RWMutex),
}

func Sniff(packetChannel <-chan gopacket.Packet) {
	parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4)
	for packet := range packetChannel {
		gotPacket(packet)
	}
}

func gotPacket(packet gopacket.Packet) {
	data := packet.Data()
	parser.DecodeLayers(data, &decoded)
	for _, layerType := range decoded {
		switch layerType {
		case layers.LayerTypeIPv4:
			go flow.Add(Endpoints{ip4.SrcIP.String(), ip4.DstIP.String()},
				Traffic(len(data)))
		case layers.LayerTypeIPv6:
			go flow.Add(Endpoints{ip6.SrcIP.String(), ip6.DstIP.String()},
				Traffic(len(data)))
		}
	}
}
