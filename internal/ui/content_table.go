package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ContentTableRow represents a single row in a content table.
type ContentTableRow struct {
	// Action is the function to be executed when this row is selected.
	Action func()
	// Data holds the textual content for each column in this row.
	Data []string
}

// ContentTableData structures the data necessary to build a content table.
type ContentTableData struct {
	// Headers are the titles for each column.
	Headers []string
	// Data is a slice of ContentTableRow, representing each row in the table.
	Data []ContentTableRow
}

// ContentTable create and return a new tview.Table widget with the provided data.
func ContentTable(data ContentTableData) *tview.Table {
	table := tview.NewTable().SetSelectable(true, false)
	columnWidth := 100 / len(data.Headers)
	table.SetBorderPadding(0, 0, 1, 1)

	for i, header := range data.Headers {
		table.SetCell(0, i, tview.NewTableCell(strings.ToUpper(header)).
			SetSelectable(false).
			SetAlign(tview.AlignLeft).
			SetMaxWidth(columnWidth).
			SetExpansion(1),
		)
	}
	for r, row := range data.Data {
		for c, col := range row.Data {
			table.SetCell(r+1, c, tview.NewTableCell(col)).SetBordersColor(tcell.ColorDodgerBlue)
		}
	}
	table.SetSelectedFunc(func(r, _ int) {
		data.Data[r-1].Action()
	})
	return table
}

// UpdateTable updates the given tview.Table with new data.
func UpdateTable(table *tview.Table, data []ContentTableRow) {
	for r, row := range data {
		for c, col := range row.Data {
			table.SetCell(r+1, c, tview.NewTableCell(col))
		}
	}
	table.SetSelectedFunc(func(r, _ int) {
		data[r-1].Action()
	})
	App.Draw()
}
