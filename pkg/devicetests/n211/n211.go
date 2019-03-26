package n211

import (
	"github.com/ExploratoryEngineering/labdevicetester/pkg/devicetests"
)

func New() *devicetests.ATDeviceTests {
	spec := devicetests.ATDeviceSpec{
		BaudRate: 9600,
		Reboot:   `AT+NRB`,
		// DisableAutoConnect: `AT+NCONFIG="AUTOCONNECT","FALSE"`,
		// EnableAutoConnect:  `AT+NCONFIG="AUTOCONNECT","TRUE"`,
		// ConfigAPN:          `AT+CGDCONT=0,"IP","%s"`,
		AutoOperatorSelection: `AT+COPS=0`,
		PSM:                   `AT+CPSMS=%d,,,"%08b","%08b"`,
		DisableEDRX:           `AT+CEDRXS=0,5`,
		CreateSocket:          `AT+NSOCR="DGRAM",17,1234,1`,
		CloseSocket:           `AT+NSOCL=%d`,
		SendUDP:               `AT+NSOST=0,"1.2.3.4",1234,2,"ABCD"`,
	}
	return devicetests.New(spec)
}
