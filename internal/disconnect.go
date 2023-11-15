package internal

import (
	"context"

	"github.com/mr-chelyshkin/NetSurf/internal/schedule"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"
)

func disconnect(ctx context.Context) {
	// ctx, cancel := context.WithCancel(ctx)

	network := make(chan string, 1)
	schedule.NetworkStatus(ctx, network)
	for {
		select {
		case _ = <-network:
			data := ui.ContentModalData{
				Action: map[string]func(){
					"cancel": func() {},
				},
			}
			ui.DrawView(ctx, "asd", ui.ContentModal(data))
		case <-ctx.Done():
			return
		}
	}
	// ctx, cancel := context.WithCancel(ctx)

	// network := make(chan string, 1)
	// schedule.NetworkStatus(ctx, network)
	// for {
	//	select {
	//	case network := <-network:
	//		modalData := ui.ContentModalData{
	//			Action: map[string]func(){
	//				"cancel": func() {
	//					cancel()
	//					return
	//				},
	//				"shit": func() {},
	//			},
	//			Text: "asd",
	//		}
	//		if network != "" {
	//			modalData.Action["disconnect"] = func() {}
	//		}

	//		ui.DrawView(ctx, "disconnect", ui.ContentModal(modalData))
	//	case <-ctx.Done():
	//		return
	//	}
	//}
}
