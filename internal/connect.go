package internal

import (
	"context"
	"fmt"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/schedule"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"
)

func connect(ctx context.Context) {
	ctx = context.WithValue(ctx, NetSurf.CtxKeyLoggerChannel, make(chan string, 1))

	view := ui.ContentTable(ui.ContentTableData{
		Headers: []string{"ssid", "freq", "quality", "level"},
	})
	ui.DrawView(ctx, view)

	ctx.Value(
		NetSurf.CtxKeyLoggerChannel,
	).(chan string) <- fmt.Sprintf(
		"scanning Wi-Fi networks every %ds",
		NetSurf.TickScanOperation,
	)

	networks := make(chan []map[string]string, 1)
	schedule.NetworkScan(ctx, networks)
	for {
		select {
		case networks := <-networks:
			data := []ui.ContentTableRow{}

			for _, network := range networks {
				data = append(data, ui.ContentTableRow{
					Action: nil,
					Data: []string{
						network["ssid"],
						network["freq"],
						network["quality"],
						network["level"],
					},
				})
			}
			ui.UpdateTable(view, data)
			ui.App.Draw()
		case <-ctx.Done():
			return
		}
	}
}
