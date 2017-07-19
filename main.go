package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/rommelc1/albionmarket-client/assemblers"
	"github.com/rommelc1/albionmarket-client/utils"
)

func main() {
	log.Print("Welcome Iron Banker")
	log.Print("You are now connected to the IronBank")
	config := utils.ClientConfig{}

	flag.StringVar(&config.DeviceName, "d", "", "Specifies the network device name. If not specified the first enumerated device will be used.")
	flag.StringVar(&config.IngestUrl, "i", "http://localhost:9000/api/marketorders/streampost/", "URL to send market data to.")
	flag.Parse()

	config.DeviceName = networkDeviceName(config.DeviceName)
	config.Username = getUser()
	if config.Username == "" {
		log.Println("Error: invalid user name")
		os.Exit(1)
	} else {
		log.Printf("Username: %s", config.Username)
	}

	log.Printf("Using the following network device: %v", config.DeviceName)
	// log.Printf("Using the following ingest: %v", config.IngestUrl)

	handle, err := pcap.OpenLive(config.DeviceName, 2048, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	var filter = "udp"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	source.NoCopy = true

	assembler := assemblers.NewMarketAssembler(config)

	log.Print("Starting to process packets...")
	for packet := range source.Packets() {
		assembler.ProcessPacket(packet)
	}
}

func networkDeviceName(deviceName string) string {
	if deviceName == "" {
		devs, err := pcap.FindAllDevs()
		if err != nil {
			log.Fatal(err)
		}
		if len(devs) == 0 {
			log.Fatal("Unable to find network device.")
		}

		if runtime.GOOS == "windows" {
			for _, device := range devs {
				// Quick and dirt hack around dealing with VirtualBox interfaces on windows
				// as one of them is often the first in the device list
				if device.Description != "Oracle" {
					return device.Name
				}
			}
		}

		return devs[0].Name
	}

	return deviceName
}

func getUser() string {
	validUsername := regexp.MustCompile(`^[A-Za-z0-9_]{3,20}$`)
	usernameLocations := [...]string{"C:\\Users\\Public\\Documents\\username.txt", "/media/username.txt"}

	for _, loc := range usernameLocations {
		if user, err := ioutil.ReadFile(loc); err == nil {
			stringUser := string(user)
			if validUsername.MatchString(stringUser) == true {
				return stringUser
			} else {
				return ""
			}
		}
	}

	return ""
}
