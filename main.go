package main

import (
	"flag"
	"fmt"
	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/phayes/freeport"
	"github.com/pkg/errors"
	"github.com/thibran/pubip"
)

var (
	w           *astilectron.Window
	l           listener
	appName     string
	builtAt     string
	isListening = true
)

func main() {
	flag.Parse()
	astilog.FlagInit()
	astilog.Debugf("Running app built at %s", builtAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName: appName,
		},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]

			// Sets host as outward IP of port 3000 as TCP
			IP, _ := pubip.NewMaster().Address()
			port, _ := freeport.GetFreePort()

			l.SetIP(IP)
			l.SetProtocol("tcp")
			l.SetPort(fmt.Sprintf("%v", port))

			go func() {
				if err := bootstrap.SendMessage(w, "ip", IP); err != nil {
					return
				}
			}()

			go func() {
				if err := bootstrap.SendMessage(w, "port", port); err != nil {
					return
				}
			}()
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: messageHandler,
			Options: &astilectron.WindowOptions{
				Title:           astilectron.PtrStr("GEM : Go Encryption Messenger"),
				BackgroundColor: astilectron.PtrStr("#efeff2"),
				Frame:           astilectron.PtrBool(false),
				Resizable:       astilectron.PtrBool(false),
				HasShadow:       astilectron.PtrBool(true),
				Fullscreenable:  astilectron.PtrBool(false),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(600),
				Width:           astilectron.PtrInt(900),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
