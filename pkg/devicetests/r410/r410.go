package r410

import (
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicetests"
)

func New() *devicetests.ATDeviceTests {
	spec := devicetests.ATDeviceSpec{
		BaudRate:              115200,
		Reboot:                `AT+CFUN=15`,
		AutoOperatorSelection: `AT+COPS=0`,
		PSM:                   `AT+CPSMS=%d,,,"%08b","%08b"`,
		DisableEDRX:           `AT+CEDRXS=0,5`,
		CreateUDPSocket:       `AT+USOCR=17,%d`,
		CreateTCPSocket:       `AT+USOCR=6,%d`,
		CloseSocket:           `AT+NSOCL=%d`,
		SendUDP:               `AT+NSOST=%d,"%v",%d,%d,"%X"`,
	}
	return devicetests.New(spec)
}
