package saran2

import (
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily"
)

func New() *devicefamily.ATdevicefamily {
	spec := devicefamily.ATDeviceSpec{
		BaudRate:        9600,
		Reboot:          `AT+NRB`,
		FirmwareVersion: `ATI9`,
		Radio:           `AT+CFUN=%v`,
		// DisableAutoConnect: `AT+NCONFIG="AUTOCONNECT","FALSE"`,
		// EnableAutoConnect:  `AT+NCONFIG="AUTOCONNECT","TRUE"`,
		ConfigAPN:                 `AT+CGDCONT=0,"IP","%s";+CGATT=1`,
		AutoOperatorSelection:     `AT+COPS=0`,
		RegistrationStatus:        `AT+CEREG?`,
		PSM:                       `AT+CPSMS=%d,,,"%08b","%08b"`,
		DisableEDRX:               `AT+CEDRXS=0,5`,
		CreateUDPSocket:           `AT+NSOCR="DGRAM",17,%d,1`,
		CloseSocket:               `AT+NSOCL=%d`,
		SendUDP:                   `AT+NSOSTF=%[1]d,"%[2]v",%[3]d,0x%03[4]x,%[5]d,"%[6]X"`,
		ReceiveUDP:                `AT+NSORF=%d,%d`,
		ReceivedMessageIndication: `+NSONMI`,
	}
	return devicefamily.New(spec)
}
