package schedule

import (
	"context"
	"os/user"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/controller"
)

// NetworkScan schedule worker.
// Scanning Wi-Fi networks (by controller.Controller) and return result to income chanel.
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

	f := func(internalCtx context.Context) {
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
		case <-internalCtx.Done():
			return
		case <-done:
			return
		}
	}
	go schedule(ctx, NetSurf.TickScanOperation, f)
}

// NetworkStatus schedule worker.
// Check current Wi-Fi connection (by controller.Controller) and return ssid to income chanel.
func NetworkStatus(ctx context.Context, c chan<- string) {
	wifi, ok := ctx.Value(NetSurf.CtxKeyWifiController).(controller.Controller)
	if !ok {
		c <- "missed controller"
		return
	}

	f := func(internalCtx context.Context) {
		done := make(chan struct{})

		go func() {
			defer close(done)

			c <- wifi.Status(ctx, nil)
		}()
		select {
		case <-internalCtx.Done():
			c <- "error"
		case <-done:
			return
		}
	}
	go schedule(ctx, NetSurf.TickCommonOperation, f)
}

// UserInfo schedule worker.
// Get current User and Permissions and return result to income channel as [2]string{}.
func UserInfo(ctx context.Context, c chan<- [2]string) {
	var (
		uid = "error"
		usr = "error"
	)

	f := func(internalCtx context.Context) {
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
		case <-internalCtx.Done():
			c <- [2]string{usr, uid}
		case <-done:
			c <- [2]string{usr, uid}
			return
		}
	}
	go schedule(ctx, NetSurf.TickCommonOperation, f)
}
