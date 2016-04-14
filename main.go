package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	dev         = `\Device\NPF_{669726D1-2173-4A7A-89DB-0843CDF048B3}`
	snapshotLen = int32(1024)
	promiscuous = false
	timeout     = 30 * time.Second
)

func main() {
	// Find all devices
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	selected := SelectDeviceFromUser(devices)

	// Open device
	handle, err := pcap.OpenLive(devices[selected].Name, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Infinity loop for sniff packets
	Sniff(packetSource.Packets())
}

func SelectDeviceFromUser(devices []pcap.Interface) (selected int) {
	ConsoleClear()

	fmt.Println(">> Please the network card to sniff packets.")
	for i, device := range devices {
		fmt.Printf("\n\t%d. Name : %s\n\t   Description : %s\n\t   IP address : %v\n",
			i+1, device.Name, device.Description, device.Addresses[1].IP)
	}
	fmt.Print("\n>> ")
	fmt.Scanf("%d", &selected)
  
  selected--

	if selected < 0 || selected > len(devices) {
		log.Panic("Invaild Selected.")
	}

	return
}
