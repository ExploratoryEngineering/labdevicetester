package devicetests

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ExploratoryEngineering/labdevicetester/pkg/serial"
)

type ATDeviceSpec struct {
	BaudRate int
	Reboot   string
	// DisableAutoConnect string
	// EnableAutoConnect  string
	// ConfigAPN          string
	AutoOperatorSelection string
	PSM                   string
	DisableEDRX           string
	CreateSocket          string
	CloseSocket           string
	SendUDP               string
}

type ATDeviceTests struct {
	s    *serial.SerialConnection
	spec ATDeviceSpec
}

func New(spec ATDeviceSpec) *ATDeviceTests {
	at := ATDeviceTests{}
	at.spec = spec
	return &at
}

func (t *ATDeviceTests) Init(s *serial.SerialConnection) {
	t.s = s
}

func (t *ATDeviceTests) BaudRate() int {
	return t.spec.BaudRate
}

func (t *ATDeviceTests) Clean() bool {

	return t.autoOperatorSelection() &&
		t.powerSaveMode(1, 255, 0) &&
		t.disableEDRX()
}

func (t *ATDeviceTests) IMEI() (int, error) {
	_, urcs, err := t.s.SendAndReceive("AT+CGSN=1")
	if err != nil {
		log.Printf("Error: %v", err)
		return 0, err
	}

	imei := strings.Split(urcs[0], ": ")[1]
	return strconv.Atoi(imei)
}

func (t *ATDeviceTests) IMSI() (int, error) {
	lines, _, err := t.s.SendAndReceive("AT+CIMI")
	if err != nil {
		log.Printf("Error: %v", err)
		return 0, err
	}
	return strconv.Atoi(lines[0])
}

func (t *ATDeviceTests) rebootModule() bool {
	log.Println("Rebooting device...")
	res, _, err := t.s.SendAndReceive(t.spec.Reboot)
	if err != nil {
		log.Printf("Error rebooting: %v", strings.Join(res, " | "))
		return false
	}
	log.Println("Rebooted OK")
	return true
}

func (t *ATDeviceTests) autoOperatorSelection() bool {
	log.Println("Auto operator selection...")
	_, _, err := t.s.SendAndReceive(t.spec.AutoOperatorSelection)
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	return true
}

// func (t *ATDeviceTests) disableAutoconnect() bool {
// 	log.Println("Disabling autoconnect...")
// 	res, _, err := t.s.SendAndReceive(t.spec.DisableAutoConnect)
// 	if err != nil {
// 		log.Printf("Error: %v (%v)", err, strings.Join(res, " | "))
// 		return false
// 	}
// 	log.Println("Autoconnect disabled")
// 	return t.rebootModule()
// }

// func (t *ATDeviceTests) configAPN() bool {
// 	log.Println("Configuring telenor.iot APN...")
// 	cmd := fmt.Sprintf(t.spec.ConfigAPN, "telenor.iot")
// 	_, _, err := t.s.SendAndReceive(cmd)
// 	if err != nil {
// 		log.Printf("Error: %v", err)
// 	}
// 	log.Println("APN configured")
// 	return true
// }

// func (t *ATDeviceTests) enableAutoconnect() bool {
// 	log.Println("Enabling autoconnect...")
// 	res, _, err := t.s.SendAndReceive(t.spec.EnableAutoConnect)
// 	if err != nil {
// 		log.Printf("Error: %v (%v)", err, strings.Join(res, " | "))
// 		return false
// 	}
// 	log.Println("Autoconnect enabled")
// 	return t.rebootModule()
// }

func (t *ATDeviceTests) powerSaveMode(enabled, tau, activeTime uint8) bool {
	log.Printf("Power save mode... %d", enabled)
	cmd := fmt.Sprintf(t.spec.PSM, enabled, tau, activeTime)
	log.Println(cmd)
	_, _, err := t.s.SendAndReceive(cmd)
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	log.Println("Power save mode configured")
	return true
}

func (t *ATDeviceTests) disableEDRX() bool {
	log.Println("Disabling eDRX...")
	_, _, err := t.s.SendAndReceive(t.spec.DisableEDRX)
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	log.Println("eDRX disabled")
	return true
}

// func (t *ATDeviceTests) sendUdpPacket(ip string, port int, data []byte) bool {
// 	log.Println("Sending tiny data packet...")
// 	response, _, err := t.s.SendAndReceive(t.spec.CreateSocket)
// 	if err != nil {
// 		log.Printf("Error creating socket: %v", err)
// 		return false
// 	}
// 	socket, err := strconv.Atoi(response[0])
// 	if err != nil {
// 		log.Printf("Error parsing socket number", err)
// 		return false
// 	}

// 	_, _, err = t.s.SendAndReceive(t.spec.SendUDP)
// 	if err != nil {
// 		log.Printf("Got error sending packet: %v", err)
// 		return false
// 	}
// 	_, _, err = t.s.SendAndReceive(fmt.Sprintf(t.spec.CloseSocket, socket))
// 	if err != nil {
// 		log.Printf("Couldn't close socket: %v", err)
// 		return false
// 	}
// 	log.Println("Successfully sent data")
// 	return true
// }
