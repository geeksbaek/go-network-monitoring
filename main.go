package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	snapshot_len int32 = 1024
	promiscuous  bool  = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
)

func main() {
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	// Open device
	handle, err = pcap.OpenLive(devices[7].Name, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case packet := <-packetSource.Packets():
			fmt.Println(packet.String())
		case <-ticker:
			fmt.Println(".")
		}
	}
}
