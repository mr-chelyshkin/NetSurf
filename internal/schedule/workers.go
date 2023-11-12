package schedule

import (
	"context"
	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/controller"
	"os/user"
)

// NetworkScan and write scan results to channel in foreground by period.
func NetworkScan(ctx context.Context, c chan<- []map[string]string) {
	output, ok := ctx.Value(NetSurf.CtxKeyLoggerChannel).(chan string)
	if !ok {
		return
	}
	wifi, ok := ctx.Value(NetSurf.CtxKeyWifiController).(controller.Controller)
	if !ok {
		output <- "missed controller"
		return
	}

	f := func(ctx context.Context) {
		done := make(chan struct{})

		go func() {
			defer close(done)

			networks := []map[string]string{}
			for _, network := range wifi.Scan(ctx, output) {
				networks = append(networks, map[string]string{
					"ssid":    network.GetSSID(),
					"freq":    network.GetFreq(),
					"level":   network.GetLevel(),
					"quality": network.GetQuality(),
				})
			}
			c <- networks
		}()
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		}
	}
	go schedule(ctx, NetSurf.TickScanOperation, f)
}

func NetworkStatus(ctx context.Context, c chan<- string) {
	wifi, ok := ctx.Value(NetSurf.CtxKeyWifiController).(controller.Controller)
	if !ok {
		c <- "missed controller"
		return
	}

	f := func(ctx context.Context) {
		done := make(chan struct{})

		go func() {
			defer close(done)

			c <- wifi.Status(ctx, nil)
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
			defer close(done)

			u, err := user.Current()
			if err == nil {
				usr = u.Username
				uid = u.Uid
			}
		}()
		select {
		case <-ctx.Done():
			c <- [2]string{usr, uid}
		case <-done:
			c <- [2]string{usr, uid}
			return
		}
	}
	go schedule(ctx, NetSurf.TickCommonOperation, f)
}
