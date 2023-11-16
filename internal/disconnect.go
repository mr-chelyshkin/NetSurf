package internal

import (
	"context"
	"fmt"
	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/controller"

	"github.com/mr-chelyshkin/NetSurf/internal/schedule"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"
)

func disconnect(ctx context.Context) {
	ctx = context.WithValue(ctx, NetSurf.CtxKeyLoggerChannel, make(chan string, 1))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	network := make(chan string, 1)
	schedule.NetworkStatus(ctx, network)

	for {
		select {
		case network := <-network:
			data := ui.ContentModalData{
				Action: map[string]func(){
					"cancel": func() {
						cancel()
						Run()
					},
				},
				Text: "no active connection",
			}
			if network != "" {
				data.Text = fmt.Sprintf("disconnect from '%s'?", network)
				data.Action["disconnect"] = func() {
					wifi, ok := ctx.Value(NetSurf.CtxKeyWifiController).(*controller.Controller)
					if !ok {
						return
					}
					wifi.Disconnect(ctx, ctx.Value(NetSurf.CtxKeyLoggerChannel).(chan string))
				}
			}
			ui.DrawView(ctx, "network_modal", ui.ContentModal(data))
		case <-ctx.Done():
			return
		}
	}
}
