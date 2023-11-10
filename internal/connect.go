package internal

import (
	"context"
	"fmt"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/schedule"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	networks := make(chan []map[string]string, 1)
	output := ctx.Value(NetSurf.CtxKeyLoggerChannel).(chan string)

	data := ui.ContentTableData{
		Headers: []string{"ssid", "freq", "quality", "level"},
		Data:    []ui.ContentTableRow{},
	}
	ui.Draw(ctx, ui.ContentTable(ctx, data))

	ctx, cancel := context.WithCancel(ctx)
	output <- fmt.Sprintf("scan wireless network every %ds", NetSurf.TickScanOperation)
	go func() {
		defer func() {
			close(networks)
			close(output)
		}()
		for {
			select {
			case networks := <-networks:
				dd := []ui.ContentTableRow{}

				for _, network := range networks {
					nn := ui.ContentTableRow{
						Action: nil,
						Data:   []string{network["ssid"], network["freq"], network["quality"], network["level"]},
					}
					dd = append(dd, nn)
				}
				data.Data = dd
				ui.App.QueueUpdateDraw(func() {
					g := ui.ContentTable(ctx, data)
					go ui.Draw(ctx, g)
				})
			case <-interrupt:
				cancel()
				return
			case <-ctx.Done():
				cancel()
				return
			}
		}
	}()
	schedule.NetworkScan(ctx, networks)
}
