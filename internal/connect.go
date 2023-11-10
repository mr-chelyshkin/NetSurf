package internal

import (
	"context"
	"fmt"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/controller"
	"github.com/mr-chelyshkin/NetSurf/internal/schedule"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"

	"github.com/rivo/tview"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	output := make(chan string, 1)
	defer close(output)

	networks := make(chan []map[string]string)
	defer close(networks)

	ctx = context.WithValue(ctx, NetSurf.CtxKeyLoggerChannel, output)
	ctx, cancel := context.WithCancel(ctx)

	view := tview.NewList()
	go ui.Draw(ctx, view)
	go func() {
		for {
			select {
			case networks := <-networks:
				output <- "networks"
				// ui.App.QueueUpdateDraw(func() {
				//	view.Clear()

				for _, network := range networks {
					output <- fmt.Sprintf("%s %s %s %s", network["ssid"], network["level"], network["quality"], network["freq"])
					//						network := network
					//						view.AddItem(network["ssid"], network["level"], '*',
					//							func() {

					//			})
				}
				//})
			case <-ctx.Done():
				return
			}
		}
	}()

	_, ok := ctx.Value(NetSurf.CtxKeyWifiController).(controller.Controller)
	if !ok {
		panic("SHHIIT")
	}
	schedule.NetworkScan(ctx, networks)

	func() {
		select {
		case <-ctx.Done():
			return
		case <-interrupt:
			cancel()
			return
		}
	}()
}
