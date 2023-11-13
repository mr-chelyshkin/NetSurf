package ui

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func mainFrame(ctx context.Context, p tview.Primitive, title string) *tview.Frame {
	content := primitive(title, p)
	header := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(info(ctx), 0, 1, false).
		AddItem(hotKeys(ctx), 0, 2, false)
	footer := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(writer(ctx), 0, 1, false)
	frame := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 6, 1, false).
		AddItem(content, 0, 3, true).
		AddItem(footer, 16, 2, false)
	f := tview.NewFrame(frame)
	f.SetBackgroundColor(tcell.ColorBlack)
	return f
}

func modalFrame(ctx context.Context, background, modal tview.Primitive, title string, w, h int) *tview.Frame {
	content := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(modal, h, 1, true).
			AddItem(nil, 0, 1, false),
			w, 1, true).
		AddItem(nil, 0, 1, false)
	container := tview.NewPages().
		AddPage("background", background, true, true).
		AddPage("modal", content, true, true)
	return mainFrame(ctx, container, title)
}
