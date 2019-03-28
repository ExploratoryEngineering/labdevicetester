package devicefamily

import (
	"github.com/ExploratoryEngineering/labdevicetester/pkg/serial"
)

type Interface interface {
	BaudRate() int
	Init(*serial.SerialConnection)
	IMEI() (int, error)
	IMSI() (int, error)
	RebootModule() bool
	SetAPN(apn string) bool
	PowerSaveMode(enabled, tau, activeTime uint8) bool
	AutoOperatorSelection() bool
	RegistrationStatus() (int, error)
	DisableEDRX() bool
	CreateSocket(protocol string, listenPort int) (int, error)
	CloseSocket(socket int) bool
	SendUDP(socket int, ip string, port int, flag SendFlag, data []byte) bool
	ReceiveUDP(socket, expectedBytes int) ([]byte, error)
}

type SendFlag int

const (
	SendFlagNone                    SendFlag = 0x000
	SendFlagHighPriority            SendFlag = 0x100
	SendFlagReleaseAfterNextMessage SendFlag = 0x200
	SendFlagReleaseAfterNextReply   SendFlag = 0x400
)
