package main

import (
	"fmt"
	"log"
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
	statistic   = make(Statistic)
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
	ticker := time.Tick(time.Second * 3)
	for {
		select {
		case packet := <-packetSource.Packets():
			gotPacket(packet)
		case <-ticker:
			fmt.Println(statistic.SortedString())
		}
	}
}

func gotPacket(packet gopacket.Packet) {
	if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
		ip, _ := ipv4Layer.(*layers.IPv4)
		traffic := len(packet.Data())
		statistic.SetTraffic(ip, uint64(traffic))
	}
}
