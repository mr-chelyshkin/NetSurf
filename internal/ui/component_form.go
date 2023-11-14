package ui

import "github.com/rivo/tview"

type ContentFormField struct {
	Label string
	Type  string
	Value string
}

type ContentFormButton struct {
	Action func()
	Label  string
}

// ContentFormData ...
type ContentFormData struct {
	Fields  []ContentFormField
	Buttons []ContentFormButton
}

// ContentForm ...
func ContentForm(data ContentFormData) *tview.Form {
	f := tview.NewForm()
	for _, field := range data.Fields {
		switch field.Type {
		case "input":
			f.AddInputField(field.Label, field.Value, 0, nil, nil)
		case "password":
			f.AddPasswordField(field.Label, field.Value, 0, '*', nil)
		}
	}
	for _, button := range data.Buttons {
		f.AddButton(button.Label, button.Action)
	}
	return f
}
