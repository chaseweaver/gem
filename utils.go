package main

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	fport = []uint16{3074, 27014, 27015, 27016, 27017,
		27018, 27019, 27020, 27021, 27022, 27023, 27024,
		27025, 27026, 27027, 27028, 27029, 27030, 27031,
		27032, 27033, 27034, 27035, 27036, 27037, 27038,
		27039, 27040, 27041, 27042, 27043, 27044, 27045,
		27046, 27047, 27048, 27049, 27050}
	password = ""
)

type (

	// Content handles message encoding
	Content struct {
		Name string `json:"name"`
		Msg  string `json:"msg"`
	}

	errorString struct {
		s string
	}
)

// New (string) error
// Returns an error that formats as the given text
func New(text string) error {
	return &errorString{text}
}

// Error() string
// Used to return a new formated error
func (e *errorString) Error() string {
	return e.s
}

// isPortOpen(string, string, string) bool
// Checks if a host+port is open, closes after successful check
func isPortOpen(protocol, host, port string) bool {
	conn, err := net.Dial(protocol, net.JoinHostPort(host, port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// outboundIP() (string, error)
// Returns public IP address
func outboundIP() (string, uint16, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", 0, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	port, err := strconv.Atoi(localAddr[1+idx:])
	if err != nil {
		return "", 0, err
	}
	return localAddr[0:idx], uint16(port), nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
