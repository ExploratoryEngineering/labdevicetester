package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicetests"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicetests/n211"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/serial"
)

var device = flag.String("device", "/dev/cu.SLAB_USBtoUART", "Serial device")
var deviceType = flag.String("type", "", "Device family type (see pkg/devicetests subfolders)")
var verbose = flag.Bool("v", false, "Verbose output")
var printIds = flag.Bool("print", false, "Print IMSI and IMEI and exit")
var serverIP = flag.String("serverip", "10.0.0.1", "IP address to the server receiving data")

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

	if !clean(tests) {
		log.Println("Clean failed")
		reportError()
		return
	}

	for {
		status, err := tests.RegistrationStatus()
		if err != nil {
			log.Println("Status failed")
			reportError()
			return
		}
		if status == 1 {
			break
		}
		log.Printf("Not connected... status %d\n", status)
		time.Sleep(1000 * time.Millisecond)
	}

	if !sendSmallPacket(tests) {
		log.Println("Sending failed")
		reportError()
		return
	}
	log.Println("Success!")
}

func clean(t devicetests.Interface) bool {
	return t.RebootModule() &&
		t.AutoOperatorSelection() &&
		t.PowerSaveMode(0, 255, 10) &&
		t.DisableEDRX()
}

func sendSmallPacket(t devicetests.Interface) bool {
	socket, err := t.CreateSocket("UDP", 1234)
	if err != nil {
		log.Printf("Error: ", err)
		reportError()
		return false
	}
	defer t.CloseSocket(socket)
	t.SendUDP(socket, *serverIP, 1234, []byte("hi"))
	return true
	return true
}
