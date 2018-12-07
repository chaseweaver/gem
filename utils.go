package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type (
	scanner struct {
		ip   string
		lock *semaphore.Weighted
	}
)

var (
	openPorts = []int{27000, 27001, 27002, 27003, 27004, 27005, 27006, 27007, 27008, 27009}
	fPort     = 1
	lPort     = 65535
	password  = ""
)

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

// outboundIP() (string, error)
// Returns public IP address
func outboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}

// portScanner()
// Starts port cycling goroutine
func portScanner() {
	IP, err := outboundIP()
	if err != nil {
		log.Fatal(err)
	}

	ps := &scanner{
		ip:   IP,
		lock: semaphore.NewWeighted(ulimit()),
	}
	ps.Start(ps.ip, 500*time.Millisecond)
}

func ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}

	s := strings.TrimSpace(string(out))
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

// scan (string, int, time.Duration) (int, bool)
// Iterates through possible ports
func scan(ip string, port int, timeout time.Duration) (int, bool) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			scan(ip, port, timeout)
		} else {
			return port, false
		}
		return 0, false
	}

	conn.Close()
	return port, true
}

// Start(string, time.Duration)
// Starts port scanning
func (s *scanner) Start(ip string, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := fPort; port <= lPort; port++ {
		s.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer s.lock.Release(1)
			defer wg.Done()
			p, ok := scan(ip, port, timeout)

			if ok {
				openPorts = append(openPorts, p)
			}

		}(port)
	}
}
