package main

import (
	"flag"
	"log"
	"time"

	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily/saran2"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily/sarar4"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/otii"
	"github.com/ExploratoryEngineering/labdevicetester/pkg/serial"
)

func main() {
	var (
		serialDevice = flag.String("device", "/dev/cu.SLAB_USBtoUART", "Serial device")
		deviceType   = flag.String("type", "", "Device family type (see pkg/devicefamily subfolders)")
		verbose      = flag.Bool("v", false, "Verbose output")
		printIds     = flag.Bool("printids", false, "Print IMSI and IMEI and exit")
		serverIP     = flag.String("serverip", "10.0.0.1", "IP address to the server receiving data")
		otiiEnabled  = flag.Bool("otii", true, "Skip Otii by setting to false")
	)
	flag.Parse()

	var device devicefamily.Interface
	switch *deviceType {
	default:
		log.Fatal("Invalid device type")
	case "n2":
		device = saran2.New()
	case "r4":
		device = sarar4.New()
	}

	otii.Init(*otiiEnabled)
	// if err := otii.EnableMainPower(); err != nil {
	// 	log.Fatal("Error enabling main power:", err)
	// }
	// defer otii.DisableMainPower()
	if err := otii.Calibrate(); err != nil {
		log.Fatal("Error calibrating:", err)
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
		return
	}

	// TODO print firmware version

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
		time.Sleep(time.Second)
	}

	recording := record(20 * time.Second)
	time.Sleep(time.Second * 30)
	for i := 0; i < 3; i++ {
		if !sendAndReceive(device, *serverIP) {
			reportError()
			return
		}
		time.Sleep(5 * time.Second)
	}
	<-recording

	// TODO print status

	// if !sendAndReceive(device, *serverIP) {
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

func clean(d devicefamily.Interface) bool {
	return d.RebootModule() &&
		d.SetRadio(devicefamily.RadioFull) &&
		d.SetAPN("telenor.iot") &&
		//d.AutoOperatorSelection() &&
		d.PowerSaveMode(1, 223, 1) &&
		d.DisableEDRX()
}

func record(duration time.Duration) chan struct{} {
	ch := make(chan struct{})
	go func() {
		otii.Record(duration)
		ch <- struct{}{}
	}()
	return ch
}

func sendSmallPacket(d devicefamily.Interface, serverIP string) bool {
	socket, err := d.CreateSocket("UDP", 1234)
	if err != nil {
		log.Println("Error: ", err)
		reportError()
		return false
	}
	defer d.CloseSocket(socket)
	d.SendUDP(socket, serverIP, 1234, devicefamily.SendFlagReleaseAfterNextMessage, []byte("hi"))
	return true
}

func sendAndReceive(d devicefamily.Interface, serverIP string) bool {
	socket, err := d.CreateSocket("UDP", 1234)
	if err != nil {
		log.Println("Error: ", err)
		reportError()
		return false
	}
	defer d.CloseSocket(socket)

	d.SendUDP(socket, serverIP, 1234, devicefamily.SendFlagReleaseAfterNextReply, []byte("echo hi"))

	d.ReceiveUDP(socket, 7)

	return true
}
