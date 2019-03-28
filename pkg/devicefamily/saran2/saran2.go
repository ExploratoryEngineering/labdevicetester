package saran2

import (
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily"
)

func New() *devicefamily.ATdevicefamily {
	spec := devicefamily.ATDeviceSpec{
		BaudRate: 9600,
		Reboot:   `AT+NRB`,
		// DisableAutoConnect: `AT+NCONFIG="AUTOCONNECT","FALSE"`,
		// EnableAutoConnect:  `AT+NCONFIG="AUTOCONNECT","TRUE"`,
		// ConfigAPN:          `AT+CGDCONT=0,"IP","%s"`,
		AutoOperatorSelection:     `AT+COPS=0`,
		RegistrationStatus:        `AT+CEREG?`,
		PSM:                       `AT+CPSMS=%d,,,"%08b","%08b"`,
		DisableEDRX:               `AT+CEDRXS=0,5`,
		CreateUDPSocket:           `AT+NSOCR="DGRAM",17,%d,1`,
		CloseSocket:               `AT+NSOCL=%d`,
		SendUDP:                   `AT+NSOST=%d,"%v",%d,%d,"%X"`,
		ReceiveUDP:                `AT+NSORF=%d,%d`,
		ReceivedMessageIndication: `+NSONMI`,
	}
	return devicefamily.New(spec)
}
