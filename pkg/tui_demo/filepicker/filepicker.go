package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// DirPickerModel represents the Bubble Tea model for directory picking.
type DirPickerModel struct {
	filepicker   filepicker.Model
	selectedFile string
	quitting     bool
	err          error
}

type clearErrorMsg struct{}

// clearErrorAfter clears the terminal after the given duration.
func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

// Init initializes Bubble Tea.
func (dirPickerModel DirPickerModel) Init() tea.Cmd {
	return dirPickerModel.filepicker.Init()
}

// Update handles messages for directory selection.
func (dirPickerModel DirPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			dirPickerModel.quitting = true
			return dirPickerModel, tea.Quit
		}
	case clearErrorMsg:
		dirPickerModel.err = nil
	}

	var cmd tea.Cmd
	dirPickerModel.filepicker, cmd = dirPickerModel.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := dirPickerModel.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		dirPickerModel.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := dirPickerModel.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		dirPickerModel.err = fmt.Errorf("'%s' is not in allowed types", path)
		dirPickerModel.selectedFile = ""
		return dirPickerModel, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return dirPickerModel, cmd
}

// View renders the TUI.
func (dirPickerModel DirPickerModel) View() string {
	if dirPickerModel.quitting {
		return ""
	}
	var stringBuilder strings.Builder
	stringBuilder.WriteString("\n  ")

	if dirPickerModel.err != nil {
		stringBuilder.WriteString(dirPickerModel.filepicker.Styles.DisabledFile.Render(dirPickerModel.err.Error()))
	} else if dirPickerModel.selectedFile == "" {
		stringBuilder.WriteString("Pick a file:")
	} else {
		stringBuilder.WriteString("Selected file: " + dirPickerModel.filepicker.Styles.Selected.Render(dirPickerModel.selectedFile))
	}

	stringBuilder.WriteString("\n\n" + dirPickerModel.filepicker.View() + "\n")
	return stringBuilder.String()
}

func main() {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
	//fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.CurrentDirectory, _ = os.Getwd()

	m := DirPickerModel{
		filepicker: fp,
	}
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(DirPickerModel)
	fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
}
