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
	go ui.DrawView(ctx, "networks", view)

	ctx.Value(
		NetSurf.CtxKeyLoggerChannel,
	).(chan string) <- fmt.Sprintf(
		"scanning Wi-Fi networks every %ds",
		NetSurf.TickScanOperation,
	)

	connForm := func(ctx context.Context, ssid string) func() {
		form := ui.ContentForm(ui.ContentFormData{
			Fields: []ui.ContentFormField{
				{
					Type:  ui.FieldInput,
					Label: "ssid",
					Value: ssid,
				},
				{
					Type:  ui.FieldInput,
					Label: "country",
					Value: "US",
				},
				{
					Type:  ui.FieldPassword,
					Label: "password",
				},
			},
		})
		form = ui.UpdateFormButtons(form, []ui.ContentFormButton{
			{
				Label:  "connect",
				Action: func() {},
			},
			{
				Label:  "cancel",
				Action: func() { go ui.DrawView(ctx, "networks", view) },
			},
		})
		return func() {
			ui.DrawModal(ctx, fmt.Sprintf("connect to %s", ssid), view, form)
		}
	}

	networks := make(chan []map[string]string, 1)
	schedule.NetworkScan(ctx, networks)
	for {
		select {
		case networks := <-networks:
			rows := []ui.ContentTableRow{}

			for _, network := range networks {
				rows = append(rows, ui.ContentTableRow{
					Action: connForm(ctx, network["ssid"]),
					Data: []string{
						network["ssid"],
						network["freq"],
						network["quality"],
						network["level"],
					},
				})
			}
			ui.UpdateTable(view, rows)
			ui.App.Draw()
		case <-ctx.Done():
			return
		}
	}
}
