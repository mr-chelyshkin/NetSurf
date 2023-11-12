package internal

import (
	"context"
	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/schedule"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	networks := make(chan []map[string]string, 1)
	output := make(chan string, 1)

	ctx, cancel := context.WithCancel(ctx)
	ctx = context.WithValue(ctx, NetSurf.CtxKeyLoggerChannel, output)

	table := ui.ContentTableData{
		Headers: []string{"ssid", "freq", "quality", "level"},
		Data:    []ui.ContentTableRow{},
	}
	view := ui.ContentTable(ctx, table)
	ui.Draw(ctx, view)

	schedule.NetworkScan(ctx, networks)
	for {
		select {
		case networks := <-networks:
			output <- "tick"

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
			ui.UpdateTable(ctx, view, data)
			ui.App.Draw()

		case <-interrupt:
			cancel()
			return
		}
	}
}
