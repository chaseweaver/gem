package main

import (
	"context"
	"encoding/json"
	"fmt"
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
		IP       string
		port     uint16
		host     string
		protocol string
		peers    []string
		builder  *network.Builder
		keys     *crypto.KeyPair
		net      *network.Network
	}
)

func (l *listener) SetIP(IP string) {
	l.IP = IP
}

func (l *listener) SetPort(p string) {
	i, _ := strconv.Atoi(p)
	l.port = uint16(i)
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
		bootstrap.SendMessage(w, "receive", fmt.Sprintf("<%s> %s", ctx.Client().ID.Address, msg.Message))
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

	l.net, _ = l.builder.Build()
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
		var peerIP string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &peerIP); err != nil {
			payload = err.Error()
			return
		}

		var port string
		for i := 0; i < len(openPorts); i++ {
			if isPortOpen(l.protocol, peerIP, openPorts[i]) {
				l.SetPort(openPorts[i])
				port = openPorts[i]
				break
			}
		}

		i, _ := strconv.Atoi(port)
		peer := network.FormatAddress(l.protocol, peerIP, uint16(i))
		initListener(l.port, l.IP, l.protocol, peer)

	case "send":
		var input string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &input); err != nil {
			payload = err.Error()
			return
		}

		ctx := network.WithSignMessage(context.Background(), true)
		l.net.Broadcast(ctx, &messages.ChatMessage{Message: input})
	}
	return
}
