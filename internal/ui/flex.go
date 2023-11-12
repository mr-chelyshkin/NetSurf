package ui

import (
	"context"
	"fmt"

	"github.com/mr-chelyshkin/NetSurf"
	"github.com/mr-chelyshkin/NetSurf/internal/schedule"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*
	Common GUI objects which are parts of GUI frame.
*/

func primitive(p tview.Primitive) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p, 0, 1, true)

	flex.SetBorder(true)
	return flex
}

func writer(ctx context.Context) *tview.Flex {
	output, ok := ctx.Value(NetSurf.CtxKeyLoggerChannel).(chan string)
	if !ok {
		return tview.NewFlex()
	}
	frame := tview.NewTextView().
		SetChangedFunc(func() { App.Draw() }).
		SetDynamicColors(true).
		ScrollToEnd()
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frame, 0, 1, false)

	go func() {
		for {
			select {
			case output := <-output:
				App.QueueUpdateDraw(func() {
					_, _ = fmt.Fprintf(frame, "%s\n", output)
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	flex.SetBorder(false)
	return flex
}

// HotKeys is a structure for storing hot key data used in a GUI helper and hot key listener.
// Expects the HotKeys structure to be present in the provided context.
type HotKeys struct {
	// Action is a function to be executed when the hot key is activated.
	Action func(ctx context.Context)
	// Description is the text describing the hot key's purpose in the GUI.
	Description string
	// Key is the key code in the format of tcell.Key.
	Key tcell.Key
}

func hotKeys(ctx context.Context) *tview.Flex {
	table := tview.NewTable()
	content, ok := ctx.Value(NetSurf.CtxKeyHotKeys).([]HotKeys)
	if ok {
		row := 0
		for _, key := range content {
			table.SetCell(row, 0, tview.NewTableCell("<"+tcell.KeyNames[key.Key]+">").SetTextColor(tcell.ColorBlue))
			table.SetCell(row, 1, tview.NewTableCell(key.Description).SetTextColor(tcell.ColorGray))
			row++
		}
	}
	App.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			for _, k := range content {
				if k.Key == event.Key() {
					k.Action(ctx)
					break
				}
			}
			return event
		},
	)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, false)

	flex.SetBorder(false)
	return flex
}

func info(ctx context.Context) *tview.Flex {
	frame := tview.NewTable()
	frame.SetCell(0, 0, tview.NewTableCell("Version:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(0, 1, tview.NewTableCell(NetSurf.Version).SetTextColor(tcell.ColorWhite))
	frame.SetCell(1, 0, tview.NewTableCell("User:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(1, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorOrangeRed))
	frame.SetCell(2, 0, tview.NewTableCell("Privileged:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(2, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorOrangeRed))
	frame.SetCell(3, 0, tview.NewTableCell("Wi-Fi network:").SetTextColor(tcell.ColorYellow))
	frame.SetCell(3, 1, tview.NewTableCell("n/a").SetTextColor(tcell.ColorOrangeRed))

	usrInfoCh := make(chan [2]string, 1)
	go func() {
		schedule.UserInfo(ctx, usrInfoCh)
		for {
			select {
			case info := <-usrInfoCh:
				// username field.
				switch info[0] {
				case "error":
					frame.GetCell(1, 1).SetText("error").SetTextColor(tcell.ColorRed)
				default:
					frame.GetCell(1, 1).SetText(info[0]).SetTextColor(tcell.ColorWhite)
				}

				// user permission field.
				switch info[1] {
				case "error":
					frame.GetCell(2, 1).SetText("error").SetTextColor(tcell.ColorRed)
				case "0":
					frame.GetCell(2, 1).SetText("yes").SetTextColor(tcell.ColorWhite)
				default:
					frame.GetCell(2, 1).SetText("run app with privileged mode").SetTextColor(tcell.ColorRed)
				}

				App.Draw()
			case <-ctx.Done():
				return
			}
		}
	}()

	networkStatusCh := make(chan string, 1)
	go func() {
		schedule.NetworkStatus(ctx, networkStatusCh)
		for {
			select {
			case network := <-networkStatusCh:
				switch network {
				case "":
					frame.GetCell(3, 1).SetText("not connected").SetTextColor(tcell.ColorOrangeRed)
				default:
					frame.GetCell(3, 1).SetText(network).SetTextColor(tcell.ColorWhite)
				}

				App.Draw()
			case <-ctx.Done():
				close(networkStatusCh)
				return
			}
		}
	}()

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(frame, 0, 1, false)
	flex.SetBorder(false)
	return flex
}
