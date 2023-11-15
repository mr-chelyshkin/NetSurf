package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ContentModalData structures the data necessary to build a content modal.
type ContentModalData struct {
	Action map[string]func()
	Text   string
}

// ContentModal create and return a new tview.Modal widget with the provided data.
func ContentModal(data ContentModalData) *tview.Modal {
	keys := []string{}
	for k := range data.Action {
		keys = append(keys, k)
	}
	return tview.NewModal().
		SetDoneFunc(func(_ int, buttonLabel string) { data.Action[buttonLabel]() }).
		SetButtonBackgroundColor(tcell.ColorDodgerBlue).
		SetBackgroundColor(tcell.ColorBlack).
		SetText(data.Text).
		AddButtons(keys)
}
