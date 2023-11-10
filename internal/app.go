package internal

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/controller"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan struct{}, 1)
	defer close(stop)

	ctx = context.WithValue(ctx, NetSurf.CtxKeyHotKeys, []ui.HotKeys{
		{
			Key:         tcell.KeyESC,
			Description: "Go to main menu",
			Action: func(ctx context.Context) {
				stop <- struct{}{}
				Run()
			},
		},
		{
			Key:         tcell.KeyCtrlC,
			Description: "Exit",
			Action: func(ctx context.Context) {
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

	output := make(chan string, 1)
	ctx = context.WithValue(ctx, NetSurf.CtxKeyLoggerChannel, output)

	go connect(ctx, stop)
	for {
		select {
		case o := <-output:
			fmt.Println(o)
		}
	}
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
