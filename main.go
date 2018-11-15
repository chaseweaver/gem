package main

import (
	"flag"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"github.com/thibran/pubip"
)

var (
	w         *astilectron.Window
	l         listener
	appName   string
	builtAt   string
	openPorts = []string{"15", "23", "8444", "8484", "3000", "3001", "3002"}
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
			go func() {
				if err := bootstrap.SendMessage(w, "ip", IP); err != nil {
					return
				}
			}()
			l.SetIP(IP)
			l.SetProtocol("tcp")

			var port string
			for i := 0; i < len(openPorts); i++ {
				if isPortOpen(l.protocol, IP, openPorts[i]) {
					l.SetPort(openPorts[i])
					port = openPorts[i]
					break
				}
			}

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
				BackgroundColor: astilectron.PtrStr("#f3f3f6"),
				Frame:           astilectron.PtrBool(false),
				Resizable:       astilectron.PtrBool(false),
				HasShadow:       astilectron.PtrBool(false),
				Fullscreenable:  astilectron.PtrBool(false),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(480),
				Width:           astilectron.PtrInt(650),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
