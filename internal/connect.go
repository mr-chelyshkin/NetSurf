package internal

import (
	"context"
	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/schedule"
	"github.com/mr-chelyshkin/NetSurf/internal/ui"
	"github.com/rivo/tview"
)

func connect(ctx context.Context, interrupt chan struct{}) {
	networks := make(chan []map[string]string, 1)
	output := ctx.Value(NetSurf.CtxKeyLoggerChannel).(chan string)

	t := tview.NewTable()
	ui.Draw(ctx, t)

	ctx, cancel := context.WithCancel(ctx)
	schedule.NetworkScan(ctx, networks)
	func() {
		for {
			select {
			case networks := <-networks:
				output <- "tick"
				for i, network := range networks {
					t.SetCell(i, 0, tview.NewTableCell(network["ssid"]))
					output <- network["ssid"]
				}
				ui.App.Draw()
			case <-interrupt:
				cancel()
			case <-ctx.Done():
				cancel()
				return
			}
		}
	}()

	//networks := make(chan []map[string]string, 1)
	//output := ctx.Value(NetSurf.CtxKeyLoggerChannel).(chan string)
	//
	//data := ui.ContentTableData{
	//	Headers: []string{"ssid", "freq", "quality", "level"},
	//	Data:    []ui.ContentTableRow{},
	//}
	//t := ui.ContentTable(ctx, data)
	//ui.Draw(ctx, t)
	//
	//ctx, cancel := context.WithCancel(ctx)
	//output <- fmt.Sprintf("scan wireless network every %ds", NetSurf.TickScanOperation)
	//
	//schedule.NetworkScan(ctx, networks)
	//func() {
	//	defer func() {
	//		close(networks)
	//		close(output)
	//	}()
	//	for {
	//		select {
	//		case networks := <-networks:
	//			//for i, network := range networks {
	//			//	t.GetCell(i, 0).SetText(network["ssid"])
	//			//	t.GetCell(i, 1).SetText(network["freq"])
	//			//	t.GetCell(i, 2).SetText(network["quality"])
	//			//	t.GetCell(i, 3).SetText(network["level"])
	//			//}
	//			//ui.App.Draw()
	//			dd := []ui.ContentTableRow{}
	//
	//			for _, network := range networks {
	//				nn := ui.ContentTableRow{
	//					Action: nil,
	//					Data:   []string{network["ssid"], network["freq"], network["quality"], network["level"]},
	//				}
	//				dd = append(dd, nn)
	//				output <- network["ssid"]
	//			}
	//			data.Data = dd
	//			ui.App.QueueUpdateDraw(func() {
	//				g := ui.ContentTable(ctx, data)
	//				go ui.Draw(ctx, g)
	//			})
	//		case <-interrupt:
	//			cancel()
	//			return
	//		case <-ctx.Done():
	//			cancel()
	//			return
	//		}
	//	}
	//}()
}
