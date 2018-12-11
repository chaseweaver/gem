package main

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	initPort = uint16(3074)
	fport    = uint16(27014)
	lport    = uint16(27050)
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
