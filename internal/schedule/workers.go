package schedule

import (
	"context"
	"os/user"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/controller"
)

func NetworkStatus(ctx context.Context, c chan<- string) {
	wifi, ok := ctx.Value(NetSurf.CtxKeyWifiController).(controller.Controller)
	if !ok {
		c <- "missed controller"
		return
	}

	f := func(ctx context.Context) {
		done := make(chan struct{})
		go func() {
			c <- wifi.Status(ctx, nil)
			close(done)
		}()
		select {
		case <-ctx.Done():
			c <- "error"
		case <-done:
			return
		}
	}
	go schedule(ctx, NetSurf.TickCommonOperation, f)
}

func UserInfo(ctx context.Context, c chan<- [2]string) {
	var (
		uid = "error"
		usr = "error"
	)

	f := func(ctx context.Context) {
		done := make(chan struct{})
		go func() {
			u, err := user.Current()
			if err == nil {
				usr = u.Username
				uid = u.Uid
			}
			close(done)
		}()
		select {
		case <-ctx.Done():
			c <- [2]string{usr, uid}
			return
		case <-done:
			c <- [2]string{usr, uid}
			return
		}
	}
	go schedule(ctx, NetSurf.TickCommonOperation, f)
}

func NetworkScan(ctx context.Context, c chan<- []map[string]string) {
	networks := []map[string]string{}
	wifi, ok := ctx.Value(NetSurf.CtxKeyWifiController).(controller.Controller)
	if !ok {
		c <- networks
		return
	}
	output, _ := ctx.Value(NetSurf.CtxKeyLoggerChannel).(chan string)

	f := func(ctx context.Context) {
		done := make(chan struct{})
		go func() {
			for _, network := range wifi.Scan(ctx, output) {
				networks = append(networks, map[string]string{
					"ssid":    network.GetSSID(),
					"quality": network.GetQuality(),
					"freq":    network.GetFreq(),
					"level":   network.GetLevel(),
				})
			}
			c <- networks
			close(done)
		}()
		select {
		case <-ctx.Done():
			c <- networks
		case <-done:
			return
		}
	}
	go schedule(ctx, 6, f)
}
