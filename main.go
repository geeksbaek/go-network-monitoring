package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	dev         = `\Device\NPF_{669726D1-2173-4A7A-89DB-0843CDF048B3}`
	snapshotLen = int32(1024)
	promiscuous = false
	timeout     = 30 * time.Second
	err         error
	handle      *pcap.Handle

	inboundTraffic  = make(map[string]uint64)
	outboundTraffic = make(map[string]uint64)
)

func main() {
	// Open device
	handle, err = pcap.OpenLive(dev, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	everyOneSec := time.Tick(1 * time.Second)
	for {
		select {
		case packet := <-packetSource.Packets():
			gotPacket(packet)
		case <-everyOneSec:
			fmt.Print("\033[H\033[2J")
			printSortedStatisticMap(inboundTraffic)
		}
	}
}

func gotPacket(packet gopacket.Packet) {
	if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
		ip, _ := ipv4Layer.(*layers.IPv4)
		traffic := len(packet.Data())
		if ip.DstIP.Equal(net.IPv4(192, 168, 0, 3)) {
			inboundTraffic[ip.SrcIP.String()] += uint64(traffic)
		} else {
			outboundTraffic[ip.DstIP.String()] += uint64(traffic)
		}
	}
}
