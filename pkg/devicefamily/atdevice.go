package devicefamily

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ExploratoryEngineering/labdevicetester/pkg/serial"
)

type ATDeviceSpec struct {
	BaudRate int
	Reboot   string
	Radio    string
	// DisableAutoConnect string
	// EnableAutoConnect  string
	FirmwareVersion           string
	ConfigAPN                 string
	AutoOperatorSelection     string
	RegistrationStatus        string
	PSM                       string
	DisableEDRX               string
	CreateUDPSocket           string
	CreateTCPSocket           string
	CloseSocket               string
	SendUDP                   string
	ReceiveUDP                string
	ReceivedMessageIndication string
}

type ATdevicefamily struct {
	s    *serial.SerialConnection
	spec ATDeviceSpec
}

func New(spec ATDeviceSpec) *ATdevicefamily {
	at := ATdevicefamily{}
	at.spec = spec
	return &at
}

func (t *ATdevicefamily) Init(s *serial.SerialConnection) {
	t.s = s
}

func (t *ATdevicefamily) BaudRate() int {
	return t.spec.BaudRate
}

func (t *ATdevicefamily) FirmwareVersion() {
	log.Printf("Firmware version")
	_, _, err := t.s.SendAndReceive(t.spec.FirmwareVersion)
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func (t *ATdevicefamily) IMEI() (int, error) {
	_, urcs, err := t.s.SendAndReceive("AT+CGSN=1")
	if err != nil {
		log.Printf("Error: %v", err)
		return 0, err
	}

	imei := strings.Split(urcs[0], ": ")[1]
	return strconv.Atoi(imei)
}

func (t *ATdevicefamily) IMSI() (int, error) {
	lines, _, err := t.s.SendAndReceive("AT+CIMI")
	if err != nil {
		log.Printf("Error: %v", err)
		return 0, err
	}
	return strconv.Atoi(lines[0])
}

func (t *ATdevicefamily) TestPowerConsumption() bool {
	// wait for connection
	// prompt for antenna attenuation change to 0 dBm
	//
	return false
}

func (t *ATdevicefamily) RebootModule() bool {
	log.Println("Rebooting device...")
	res, _, err := t.s.SendAndReceive(t.spec.Reboot)
	if err != nil {
		log.Printf("Error rebooting: %v", strings.Join(res, " | "))
		return false
	}
	log.Println("Rebooted OK")
	return true
}

func (t *ATdevicefamily) SetAPN(apn string) bool {
	log.Printf("Set APN to %s...", apn)
	_, _, err := t.s.SendAndReceive(fmt.Sprintf(t.spec.ConfigAPN, apn))
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	return true
}

func (t *ATdevicefamily) SetRadio(fun RadioFunctionality) bool {
	log.Println("Radio functionality")
	radioFun := ""
	switch fun {
	case RadioOff:
		radioFun = "0"
	case RadioFull:
		radioFun = "1"
	default:
		log.Fatalln("Radio functionality not implemented")
		return false
	}
	cmd := fmt.Sprintf(t.spec.Radio, radioFun)
	_, _, err := t.s.SendAndReceive(cmd)
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	return true
}

func (t *ATdevicefamily) AutoOperatorSelection() bool {
	log.Println("Auto operator selection...")
	_, _, err := t.s.SendAndReceive(t.spec.AutoOperatorSelection)
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	return true
}

func (t *ATdevicefamily) RegistrationStatus() (int, error) {
	log.Println("Registration status...")
	_, urcs, err := t.s.SendAndReceive(t.spec.RegistrationStatus)
	if err != nil {
		log.Printf("Error: %v", err)
		return 0, err
	}
	for _, urc := range urcs {
		if strings.Index(urc, "+CEREG") != 0 {
			continue
		}
		status := strings.Split(urc, ",")[1]
		return strconv.Atoi(status)
	}
	log.Println("Error: +CEREG response not found")
	return 0, errors.New("+CEREG response not found")
}

// func (t *ATdevicefamily) disableAutoconnect() bool {
// 	log.Println("Disabling autoconnect...")
// 	res, _, err := t.s.SendAndReceive(t.spec.DisableAutoConnect)
// 	if err != nil {
// 		log.Printf("Error: %v (%v)", err, strings.Join(res, " | "))
// 		return false
// 	}
// 	log.Println("Autoconnect disabled")
// 	return t.rebootModule()
// }

// func (t *ATdevicefamily) enableAutoconnect() bool {
// 	log.Println("Enabling autoconnect...")
// 	res, _, err := t.s.SendAndReceive(t.spec.EnableAutoConnect)
// 	if err != nil {
// 		log.Printf("Error: %v (%v)", err, strings.Join(res, " | "))
// 		return false
// 	}
// 	log.Println("Autoconnect enabled")
// 	return t.rebootModule()
// }

func (t *ATdevicefamily) PowerSaveMode(enabled, tau, activeTime uint8) bool {
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

func (t *ATdevicefamily) DisableEDRX() bool {
	log.Println("Disabling eDRX...")
	_, _, err := t.s.SendAndReceive(t.spec.DisableEDRX)
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}
	log.Println("eDRX disabled")
	return true
}

func (t *ATdevicefamily) CreateSocket(protocol string, listenPort int) (int, error) {
	log.Printf("Create socket")

	var cmd string
	switch protocol {
	case "UDP":
		if t.spec.CreateUDPSocket == "" {
			log.Fatalf("Error: device does not implement UDP socket")
		}
		cmd = fmt.Sprintf(t.spec.CreateUDPSocket, listenPort)
	case "TCP":
		if t.spec.CreateTCPSocket == "" {
			log.Fatalf("Error: device does not implement TCP socket")
		}
		cmd = fmt.Sprintf(t.spec.CreateTCPSocket, listenPort)
	default:
		log.Fatal("protocol not implemented")
	}

	lines, urcs, err := t.s.SendAndReceive(cmd)
	if err != nil {
		log.Printf("Error creating socket: %v", err)
		return 0, err
	}

	var socket int
	if len(lines) > 0 {
		socket, err = strconv.Atoi(lines[0])
		if err != nil {
			log.Printf("Error parsing socket number", err)
			return 0, err
		}
	} else if len(urcs) > 0 && strings.HasPrefix(urcs[0], "+USOCR") {
		socket, err = strconv.Atoi(urcs[0][8:])
		if err != nil {
			log.Printf("Error parsing +USOCR socket number", err)
			return 0, err
		}
	}

	return socket, nil
}

func (t *ATdevicefamily) CloseSocket(socket int) bool {
	_, _, err := t.s.SendAndReceive(fmt.Sprintf(t.spec.CloseSocket, socket))
	if err != nil {
		log.Printf("Couldn't close socket: %v", err)
		return false
	}
	return true
}

func (t *ATdevicefamily) SendUDP(socket int, ip string, port int, flag SendFlag, data []byte) bool {
	log.Println("Sending UDP packet...")

	cmd := fmt.Sprintf(t.spec.SendUDP, socket, ip, port, flag, len(data), data)
	_, _, err := t.s.SendAndReceive(cmd)
	if err != nil {
		log.Printf("Error sending packet: %v", err)
		return false
	}

	log.Println("Successfully sent data")
	return true
}

func (t *ATdevicefamily) ReceiveUDP(socket, expectedBytes int) ([]byte, error) {
	log.Println("Receiving UDP Packet...")

	if t.spec.ReceivedMessageIndication != "" {
		line, err := t.s.WaitForURC(t.spec.ReceivedMessageIndication)
		if err != nil {
			log.Printf("Error receive URC: %v", err)
			return nil, err
		}
		line = strings.TrimPrefix(line, t.spec.ReceivedMessageIndication)
		log.Println(line)
	}

	cmd := fmt.Sprintf(t.spec.ReceiveUDP, socket, expectedBytes)
	_, _, err := t.s.SendAndReceive(cmd)
	if err != nil {
		log.Printf("Error receiving UDP: %v", err)
		return nil, err
	}
	return []byte{}, nil
}
