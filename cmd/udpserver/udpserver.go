package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	startTime := time.Now().Format("2006-01-02 15:04:05")
	logFile, err := os.OpenFile("udpserver "+startTime+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	mw := io.MultiWriter(os.Stdout, logFile)
	if err != nil {
		log.Fatal("Unable to open log file:", err)
	}
	log.SetOutput(mw)
	log.SetFlags(log.Ltime)

	udpAddr := ":1234"

	addr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		panic(err)
	}

	serverConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	log.Println("Starting UDP server")
	log.Printf("Listening on %s\n", udpAddr)

	for {
		buf := make([]byte, 4096)
		n, fromAddr, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		if n > 0 {
			received := string(buf[0:n])
			log.Printf("Got %d bytes from %s: %s\n", n, fromAddr, received)

			if strings.HasPrefix(received, "echo ") {
				resp := buf[5:n]
				log.Printf("Echoing %q to %v", resp, fromAddr)
				serverConn.WriteToUDP(resp, fromAddr)
			}
		}
	}
}
