package ui

import (
	"github.com/rivo/tview"
)

/*
	Main GUI objects which are parts of GUI frame.
*/

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

// ContentTable creates and returns a new tview.Table widget populated with the provided data.
func ContentTable(data ContentTableData) *tview.Table {
	content := tview.NewTable().SetSelectable(true, false)
	for i, header := range data.Headers {
		content.SetCell(0, i, tview.NewTableCell(header).SetAlign(tview.AlignCenter).SetSelectable(false))
	}
	for r, row := range data.Data {
		for c, col := range row.Data {
			content.SetCell(r+1, c, tview.NewTableCell(col))
		}
	}
	content.SetSelectedFunc(func(r, c int) {
		data.Data[r-1].Action()
	})
	return content
}

// UpdateTable updates the given tview.Table with new data.
func UpdateTable(table *tview.Table, data []ContentTableRow) {
	for r, row := range data {
		for c, col := range row.Data {
			table.SetCell(r+1, c, tview.NewTableCell(col))
		}
	}
	table.SetSelectedFunc(func(r, c int) {
		data[r-1].Action()
	})
	App.Draw()
}
