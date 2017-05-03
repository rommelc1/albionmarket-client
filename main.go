package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"github.com/google/gopacket/layers"
	"reflect"
	"encoding/binary"
)

func main() {
	log.Print("Starting things...")
	// This set of bytes indicates the start of market JSON
	match := []byte{243, 3, 1, 0, 0, 42, 0, 2, 0, 121}

	handle, err := pcap.OpenLive("en0", 1024, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	var filter string = "udp and src host 158.85.26.38"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	source.NoCopy = true

	log.Println("Starting to read packets...")
	for packet := range source.Packets() {
		udp := packet.TransportLayer().(*layers.UDP)

		if udp.Length < 54 {
			continue
		}

		marketHeader := udp.Payload[44:54]

		if reflect.DeepEqual(marketHeader, match) {
			// TODO: Create assembler thing
			// pass every packet to assembler
			// assembler puts pieces together
			// assembler detects new market response
			// starts parsing market data
			// when complete assembler sends market data off and resets assembler

			// Prototype stuff here
			numItems := udp.Payload[54:56]
			fmt.Println(numItems)

			firstEntryLength := udp.Payload[57:59]
			fmt.Println(binary.BigEndian.Uint16(firstEntryLength))
			fmt.Println(string(udp.Payload[59:59+binary.BigEndian.Uint16(firstEntryLength)]))
		}
	}
}