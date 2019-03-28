package sarar4

import (
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily"
)

func New() *devicefamily.ATdevicefamily {
	spec := devicefamily.ATDeviceSpec{
		BaudRate:              115200,
		Reboot:                `AT+CFUN=15`,
		ConfigAPN:             `AT+CGDCONT=0,"IP","%s"`,
		AutoOperatorSelection: `AT+COPS=0`,
		PSM:                   `AT+CPSMS=%d,,,"%08b","%08b"`,
		DisableEDRX:           `AT+CEDRXS=0,5`,
		CreateUDPSocket:       `AT+USOCR=17,%d`,
		CreateTCPSocket:       `AT+USOCR=6,%d`,
		CloseSocket:           `AT+NSOCL=%d`,
		SendUDP:               `AT+NSOST=%d,"%v",%d,%d,"%X"`,
	}
	return devicefamily.New(spec)
}
