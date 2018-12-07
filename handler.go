package main

import (
	"context"
	"encoding/json"
	"github.com/perlin-network/noise/crypto"
	"strconv"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/perlin-network/noise/crypto/ed25519"
	"github.com/perlin-network/noise/examples/chat/messages"
	"github.com/perlin-network/noise/network"
	"github.com/perlin-network/noise/network/discovery"
	"github.com/perlin-network/noise/types/opcode"
)

type (
	chatPlugin struct {
		*network.Plugin
	}

	listener struct {
		ip       string
		port     uint16
		host     string
		protocol string
		peers    []string
		builder  *network.Builder
		keys     *crypto.KeyPair
		net      *network.Network
	}
)

func (l *listener) SetIP(ip string) {
	l.ip = ip
}

func (l *listener) SetPort(p int) {
	l.port = uint16(p)
}

func (l *listener) SetHost(h string) {
	l.host = h
}

func (l *listener) SetProtocol(p string) {
	l.protocol = p
}

func (l *listener) AddPeer(p string) {
	l.peers = append(l.peers, p)
}

func (l *listener) SetPeers(p []string) {
	l.peers = p
}

func (l *listener) SetBuilder(b *network.Builder) {
	l.builder = b
}

func (l *listener) SetKeys(k *crypto.KeyPair) {
	l.keys = k
}

func (state *chatPlugin) Receive(ctx *network.PluginContext) error {
	switch msg := ctx.Message().(type) {
	case *messages.ChatMessage:
		// Hash key from recieved password
		key := []byte(createHash(password))
		var dmsg string
		dmsg, err := decryptMessage(msg.Message, key)
		if err != nil {
			return err
		}
		bootstrap.SendMessage(w, "receive", dmsg)
	}
	return nil
}

// initListener
// Initializes a listener on set Protocol://IP:Port and initializes peers
func initListener(port uint16, host, protocol string, peers ...string) {

	// Generate Public and Private key pair
	l.keys = ed25519.RandomKeyPair()

	opcode.RegisterMessageType(opcode.Opcode(1000), &messages.ChatMessage{})
	l.builder = network.NewBuilder()
	l.builder.SetKeys(l.keys)
	l.builder.SetAddress(network.FormatAddress(protocol, host, port))

	// Register peer discovery plugin, custom chat plugin.
	l.builder.AddPlugin(new(discovery.Plugin))
	l.builder.AddPlugin(new(chatPlugin))

	var err error
	l.net, err = l.builder.Build()
	if err != nil {
		return
	}
	go l.net.Listen()

	// Initialize peer group
	if len(peers) > 0 {
		l.net.Bootstrap(peers...)
	}
}

// messageHandler(astilectron.Window, bootstrap.MessageIn) (interface{}, error)
// Handles returned messages from JS client, returns errors, success
func messageHandler(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "close":
		w.Close()

	case "connect":
		var peer []string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &peer); err != nil {
			payload = err.Error()
			return
		}

		var ip, port string
		if len(peer) == 2 {
			ip = peer[0]
			port = peer[1]
		}

		p, _ := strconv.Atoi(port)
		peerAddress := network.FormatAddress(l.protocol, ip, uint16(p))

		if len(peerAddress) > 0 {
			l.net.Bootstrap(peerAddress)
		}

	case "send":
		var input []string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &input); err != nil {
			payload = err.Error()
			return
		}

		var msg, pwd string
		if len(input) == 2 {
			msg = input[0]
			pwd = input[1]
		}

		password = pwd
		key := []byte(createHash(pwd))
		var emsg string
		if emsg, err = encryptMessage(msg, key); err != nil {
			payload = err.Error()
		}

		ctx := network.WithSignMessage(context.Background(), true)
		l.net.Broadcast(ctx, &messages.ChatMessage{Message: emsg})
	}
	return
}
