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
	stop := make(chan struct{}, 1)
	output := make(chan string, 1)

	ctx := context.WithValue(context.Background(), NetSurf.CtxKeyHotKeys, []ui.HotKeys{
		{
			Key:         tcell.KeyEsc,
			Description: "Go to main menu",
			Action: func(context.Context) {
				stop <- struct{}{}
				Run()
			},
		},
		{
			Key:         tcell.KeyCtrlC,
			Description: "Exit",
			Action: func(context.Context) {
				stop <- struct{}{}
				ui.App.Stop()
				os.Exit(0)
			},
		},
	})
	ctx = context.WithValue(ctx, NetSurf.CtxKeyWifiController, controller.New(
		controller.WithScanSkipEmptySSIDs(),
		controller.WithScanSortByLevel(),
	))
	ctx = context.WithValue(ctx, NetSurf.CtxKeyLoggerChannel, output)

	view := ui.ContentTable(ctx, ui.ContentTableData{
		Headers: []string{"connect", "scan and connect to sifi network"},
		Data: []ui.ContentTableRow{
			{
				Action: func() { go connect(ctx, stop) },
				Data:   []string{"connect", "scan and connect to wi-fi"},
			},
			{
				Action: func() {},
				Data:   []string{"disconnect", "interrupt current wifi connection"},
			},
		},
	})
	return ui.Start(ctx, view)
}

func Run1() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, NetSurf.CtxKeyWifiController, controller.New(
		controller.WithScanSkipEmptySSIDs(),
		controller.WithScanSortByLevel(),
	))

	output := make(chan string, 1)
	ctx = context.WithValue(ctx, NetSurf.CtxKeyLoggerChannel, output)

	//view := ui.ContentTable(ctx, ui.ContentTableData{
	//	Headers: []string{"action", "description"},
	//	Data: []ui.ContentTableRow{
	//		{
	//			Action: func() { go connect(ctx, stop) },
	//			Data:   []string{"connect", "scan and connect to wifi network"},
	//		},
	//		{
	//			Action: func() {},
	//			Data:   []string{"disconnect", "interrupt current wifi connection"},
	//		},
	//	},
	//})
	//return ui.Start(ctx, view)
	return nil
}
