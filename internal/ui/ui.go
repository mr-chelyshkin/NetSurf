package ui

import (
	"context"

	"github.com/rivo/tview"
)

// App is a main application view object.
var App = tview.NewApplication()

// StartView run console GUI.
// Execute only in main process, extend mainFrame GUI layout.
func StartView(ctx context.Context, title string, p tview.Primitive) error {
	return setFrame(mainFrame(ctx, p, title)).Run()
}

// DrawView run new GUI view with income tview.Primitive as main object.
// Execute on cli commands, extend mainFrame GUI layout.
func DrawView(ctx context.Context, title string, p tview.Primitive) {
	setFrame(mainFrame(ctx, p, title)).Draw()
}

func setFrame(frame *tview.Frame) *tview.Application {
	return App.SetRoot(frame, true).SetFocus(frame)
}
