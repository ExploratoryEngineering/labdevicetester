package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily/saran2"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/serial"
)

func main() {
	var (
		serialDevice = flag.String("device", "/dev/cu.SLAB_USBtoUART", "Serial device")
		deviceType   = flag.String("type", "", "Device family type (see pkg/devicefamily subfolders)")
		verbose      = flag.Bool("v", false, "Verbose output")
		printIds     = flag.Bool("print", false, "Print IMSI and IMEI and exit")
		serverIP     = flag.String("serverip", "10.0.0.1", "IP address to the server receiving data")
	)
	flag.Parse()

	var device devicefamily.Interface
	switch *deviceType {
	default:
		log.Fatal("Invalid device type")
	case "n211":
		device = saran2.New()
	}

	s, err := serial.NewSerialConnection(*serialDevice, device.BaudRate(), *verbose)
	if err != nil {
		log.Println("Unable to open serial port:", err)
		return
	}
	defer s.Close()

	device.Init(s)

	if !checkSerial(s) {
		reportError()
		return
	}

	if *printIds {
		imsi, err := device.IMSI()
		if err != nil {
			log.Println("Error: ", err)
		}
		imei, err := device.IMEI()
		if err != nil {
			log.Println("Error: ", err)
		}

		log.Println("IMSI:", imsi)
		log.Println("IMEI:", imei)
		os.Exit(0)
	}

	if !clean(device) {
		log.Println("Clean failed")
		reportError()
		return
	}

	for {
		status, err := device.RegistrationStatus()
		if err != nil {
			log.Println("Status failed")
			reportError()
			return
		}
		if status == 1 {
			break
		}
		log.Println("Not connected... status:", status)
		time.Sleep(1000 * time.Millisecond)
	}

	if !sendSmallPacket(device, *serverIP) {
		log.Println("Sending failed")
		reportError()
		return
	}

	// if !sendAndReceive(device) {
	// 	log.Println("Send and receive failed")
	// 	reportError()
	// 	return
	// }
	log.Println("Success!")
}

func checkSerial(s *serial.SerialConnection) bool {
	log.Println("Testing serial device...")
	_, _, err := s.SendAndReceive("AT")
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	log.Println("Device responds OK")
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

func clean(t devicefamily.Interface) bool {
	return t.RebootModule() &&
		t.AutoOperatorSelection() &&
		t.PowerSaveMode(0, 255, 10) &&
		t.DisableEDRX()
}

func sendSmallPacket(t devicefamily.Interface, serverIP string) bool {
	socket, err := t.CreateSocket("UDP", 1234)
	if err != nil {
		log.Println("Error: ", err)
		reportError()
		return false
	}
	defer t.CloseSocket(socket)
	t.SendUDP(socket, serverIP, 1234, []byte("hi"))
	return true
}

// func sendAndReceive(t devicefamily.Interface) bool {
// 	socket, err := t.CreateSocket("UDP", 1234)
// 	if err != nil {
// 		log.Println("Error: ", err)
// 		reportError()
// 		return false
// 	}
// 	defer t.CloseSocket(socket)

// 	t.SendUDP(socket, *serverIP, 1234, []byte("echo hi"))

// 	t.ReceiveUDP(socket, 7)

// 	return true
// }
