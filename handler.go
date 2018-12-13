package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/perlin-network/noise/crypto"
	"io/ioutil"
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

func (l *listener) SetPort(p uint16) {
	l.port = p
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

		arr := []byte(msg.Message)
		var dat map[string]interface{}
		json.Unmarshal([]byte(arr), &dat)

		dmsg, err := decryptMessage(dat["msg"].(string), key)
		if err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
			return err
		}

		dname, err := decryptMessage(dat["name"].(string), key)
		if err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
			return err
		}

		dat["msg"] = dmsg
		dat["name"] = dname

		mdat, _ := json.Marshal(dat)

		bootstrap.SendMessage(w, "receive", string(mdat[:]))
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
		bootstrap.SendMessage(w, "error", err.Error())
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
			bootstrap.SendMessage(w, "error", err.Error())
			return
		}

		var ip, port, pwd string
		if len(peer) == 3 {
			ip = peer[0]
			port = peer[1]
			pwd = peer[2]
		}

		password = pwd

		p, _ := strconv.Atoi(port)
		peerAddress := network.FormatAddress(l.protocol, ip, uint16(p))

		if len(peerAddress) > 0 {
			l.net.Bootstrap(peerAddress)
		}

	case "change-port":
		var port string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &port); err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
			return
		}

		var value uint64
		if value, err = strconv.ParseUint(port, 10, 16); err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
			return
		}

		l.net.Close()

		bootstrap.SendMessage(w, "warn", fmt.Sprintf("Shutting down server: %v://%v:%v", l.protocol, l.ip, l.port))

		l.SetPort(uint16(value))

		bootstrap.SendMessage(w, "success", fmt.Sprintf("Listening for peers: %v://%v:%v", l.protocol, l.ip, l.port))
		initListener(l.port, l.ip, l.protocol, "")

	case "save":
		var input []string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &input); err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
			return
		}

		var filename, msgs string
		if len(input) == 2 {
			filename = input[0]
			msgs = input[1]
		}

		data := []byte(msgs)
		if err = ioutil.WriteFile(filename, data, 0644); err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
			return
		}

	case "send":
		var input []string

		// Unmarshal JSON string
		if err = json.Unmarshal([]byte(m.Payload), &input); err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
			return
		}

		var msg, pwd, nme string
		if len(input) == 3 {
			msg = input[0]
			pwd = input[1]
			nme = input[2]
		}

		password = pwd
		key := []byte(createHash(pwd))

		var emsg string
		if emsg, err = encryptMessage(msg, key); err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
		}

		var ename string
		if ename, err = encryptMessage(nme, key); err != nil {
			bootstrap.SendMessage(w, "error", err.Error())
		}

		content := Content{
			Name: ename,
			Msg:  emsg,
		}

		mcontent, _ := json.Marshal(content)
		ctx := network.WithSignMessage(context.Background(), true)
		l.net.Broadcast(ctx, &messages.ChatMessage{Message: string(mcontent[:])})
	}
	return
}
