package serial

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// SerialConnection is a serial connection
type SerialConnection struct {
	serialPort *serial.Port
	scanner    *bufio.Scanner
	verbose    bool
}

// NewSerialConnection creates a new SerialConnection
func NewSerialConnection(device string, baud int, verbose bool) (*SerialConnection, error) {
	c := &serial.Config{Name: device, Baud: baud, ReadTimeout: time.Second * 30}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}

	// Wrap serial connection in scanner
	scanner := bufio.NewScanner(s)
	scanner.Split(scanCRLF)

	return &SerialConnection{
		serialPort: s,
		scanner:    scanner,
		verbose:    verbose,
	}, nil
}

// SendAndReceive sends and recieves data, both regular commands and URCs
func (s *SerialConnection) SendAndReceive(cmd string) ([]string, []string, error) {
	if s.verbose {
		log.Printf("--> %s", cmd)
	}

	_, err := s.serialPort.Write([]byte(cmd + "\r\n"))
	if err != nil {
		return nil, nil, err
	}

	return s.scanResponse(s.scanner)
}

func (s *SerialConnection) WaitForURC(urc string) (string, error) {
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if s.verbose && line != "" {
			log.Printf("<-- %s", line)
		}

		if strings.HasPrefix(line, urc) {
			return line, nil
		}
	}
	return "", fmt.Errorf("Error: serial closed")
}

// Close closes the serial connection
func (s *SerialConnection) Close() {
	s.serialPort.Close()
}

func (s *SerialConnection) splitURCResponse(cmds []string, err error) ([]string, []string, error) {
	var urcs []string
	var data []string
	for _, v := range cmds {
		if strings.HasPrefix(v, "+") {
			urcs = append(urcs, v)
			continue
		}
		if strings.TrimSpace(v) == "" {
			continue
		}
		data = append(data, v)
	}
	time.Sleep(20 * time.Millisecond)
	return data, urcs, err
}

func (s *SerialConnection) scanResponse(scanner *bufio.Scanner) ([]string, []string, error) {
	var data []string

	for s.scanner.Scan() {
		line := scanner.Text()
		if s.verbose && line != "" {
			log.Printf("<-- %s", line)
		}

		if line == "OK" {
			return s.splitURCResponse(data[1:], nil)
		}

		if strings.Contains(line, "ERROR") {
			return s.splitURCResponse(data, fmt.Errorf("ERROR: '%v'", line))
		}

		if line == "ABORT" {
			return s.splitURCResponse(data, fmt.Errorf("ABORT: '%v'", line))
		}
		data = append(data, line)
	}

	return s.splitURCResponse(data, fmt.Errorf("Invalid response: '%v'", data))
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func scanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
