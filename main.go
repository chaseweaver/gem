package main

import (
	"flag"
	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
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

			resp, rerr := http.Get("https://api.ipify.org")
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			ip := string(body)
			if rerr != nil {
				ip, _, _ = outboundIP()
			}
			defer resp.Body.Close()

			var port uint16
			for _, v := range fport {
				if isPortOpen("tcp", ip, string(v)) {
					port = v
				}
			}

			if port == 0 {
				_, port, _ = outboundIP()
			}

			l.SetIP(ip)
			l.SetProtocol("tcp")
			l.SetPort(port)
			initListener(l.port, l.ip, l.protocol, "")

			go func() {
				if err := bootstrap.SendMessage(w, "ip", ip); err != nil {
					bootstrap.SendMessage(w, "error", err.Error())
					return
				}

				if err := bootstrap.SendMessage(w, "port", port); err != nil {
					bootstrap.SendMessage(w, "error", err.Error())
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
				Height:          astilectron.PtrInt(632),
				Width:           astilectron.PtrInt(900),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
