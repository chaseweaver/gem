package main

import "net"

// errorString struct
// Used for implementation of an error message
type errorString struct {
	s string
}

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
	conn, err := net.Listen(protocol, net.JoinHostPort(host, port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
