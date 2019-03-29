package sarar4

import (
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicefamily"
)

func New() *devicefamily.ATdevicefamily {
	spec := devicefamily.ATDeviceSpec{
		BaudRate:              115200,
		Reboot:                `AT+COPS=2;+URAT=8;+CFUN=15`,
		FirmwareVersion:       `ATI9`,
		ConfigAPN:             `AT+CGDCONT=1,"IP","%s"`,
		Radio:                 `ATE0+CFUN=%v`,
		AutoOperatorSelection: `AT+COPS=0`,
		RegistrationStatus:    `AT+CEREG?`,
		PSM:                   `AT+CPSMS=%d,,,"%08b","%08b"`,
		DisableEDRX:           `AT+CEDRXS=0,5`,
		CreateUDPSocket:       `AT+USOCR=17,%d`,
		CreateTCPSocket:       `AT+USOCR=6,%d`,
		CloseSocket:           `AT+USOCL=%d`,
		SendUDP:               `AT+USOST=%[1]d,"%[2]v",%[3]d,%[5]d,"%[6]s"`,
		ReceiveUDP:            `AT+USORF=%d,%d`,
	}
	return devicefamily.New(spec)
}
