package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type field int

const (
	FieldInput field = iota
	FieldPassword
)

// ContentFormField describe field in form.
type ContentFormField struct {
	// Label of the field.
	Label string
	// Value of the field.
	Value string
	// Type of the field.
	Type field
}

// ContentFormButton describe button in form.
type ContentFormButton struct {
	// Action func for button.
	Action func()
	// Label of the button.
	Label string
}

// ContentFormData describe form.
type ContentFormData struct {
	Fields []ContentFormField
}

// ContentForm create and return a new tview.Form widget with the provided data.
func ContentForm(data ContentFormData) *tview.Form {
	form := tview.NewForm().
		SetButtonBackgroundColor(tcell.ColorDodgerBlue).
		SetFieldBackgroundColor(tcell.ColorDodgerBlue)
	for _, field := range data.Fields {
		switch field.Type {
		case FieldInput:
			form.AddInputField(field.Label, field.Value, 0, nil, nil)
		case FieldPassword:
			form.AddPasswordField(field.Label, field.Value, 0, '*', nil)
		default:
		}
	}
	return form
}

// UpdateFormButtons add cations to the form.
func UpdateFormButtons(form *tview.Form, buttons []ContentFormButton) *tview.Form {
	for _, button := range buttons {
		form.AddButton(button.Label, button.Action)
	}
	return form
}
