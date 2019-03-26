package devicetests

import (
	"github.com/ExploratoryEngineering/labdevicetester/pkg/serial"
)

type Interface interface {
	BaudRate() int
	Init(*serial.SerialConnection)
	Clean() bool
	IMEI() (int, error)
	IMSI() (int, error)
}
