package main

import (
	"flag"
	"log"
	"os"

	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicetests"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicetests/n211"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/serial"
)

var device = flag.String("device", "/dev/cu.SLAB_USBtoUART", "Serial device")
var deviceType = flag.String("type", "", "Device family type (see pkg/devicetests subfolders)")
var verbose = flag.Bool("v", false, "Verbose output")
var printIds = flag.Bool("print", false, "Print IMSI and IMEI and exit")

func checkSerial(s *serial.SerialConnection) bool {
	log.Println("Testing serial device...")
	_, _, err := s.SendAndReceive("AT")
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	log.Printf("Device responds OK")
	return true
}

func reportError() {
	log.Println()
	log.Println("=======================================")
	log.Println("X X X X X X X X X X X X X X X X X X X X")
	log.Println("           o h    c r a p              ")
	log.Println()
	log.Println("            Test failed.")
	log.Println("X X X X X X X X X X X X X X X X X X X X")
	log.Println("=======================================")
}

func main() {
	var tests devicetests.Interface

	flag.Parse()

	switch *deviceType {
	default:
		log.Fatal("Invalid device type")
	case "n211":
		tests = n211.New()
	}

	s, err := serial.NewSerialConnection(*device, tests.BaudRate(), *verbose)
	if err != nil {
		log.Printf("Unable to open serial port: %v", err)
		return
	}
	defer s.Close()

	tests.Init(s)

	if !checkSerial(s) {
		reportError()
		return
	}

	if *printIds {
		imsi, err := tests.IMSI()
		if err != nil {
			log.Printf("Error: ", err)
		}
		imei, err := tests.IMEI()
		if err != nil {
			log.Printf("Error: ", err)
		}

		log.Printf("IMSI: %d\n", imsi)
		log.Printf("IMEI: %d\n", imei)
		os.Exit(0)
	}

	if !tests.Clean() {
		log.Fatal("Clean failed")
	}
	log.Println("Success!")
}
