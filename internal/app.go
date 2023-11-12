package internal

import (
	"context"
	"os"

	"github.com/gdamore/tcell/v2"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/controller"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, NetSurf.CtxKeyHotKeys, []ui.HotKeys{
		{
			Key:         tcell.KeyEsc,
			Description: "Go to main menu",
			Action: func(context.Context) {
				cancel()
				_ = Run()
			},
		},
		{
			Key:         tcell.KeyCtrlC,
			Description: "Exit",
			Action: func(context.Context) {
				cancel()
				ui.App.Stop()
				os.Exit(0)
			},
		},
	})
	ctx = context.WithValue(ctx, NetSurf.CtxKeyWifiController, controller.New(
		controller.WithScanSkipEmptySSIDs(),
		controller.WithScanSortByLevel(),
	))

	view := ui.ContentTable(ui.ContentTableData{
		Headers: []string{"connect", "description"},
		Data: []ui.ContentTableRow{
			{
				Action: func() { go connect(ctx) },
				Data:   []string{"connect", "scan and connect to Wi-Fi network"},
			},
			{
				Action: func() {},
				Data:   []string{"disconnect", "interrupt current Wi-Fi connection"},
			},
		},
	})
	return ui.StartView(ctx, view)
}
