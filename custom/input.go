package custom

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func SearchBar() (input textinput.Model) {
	input = textinput.New()
	input.Prompt ="🐨>"
	// input.Placeholder = "container"
	// ti.Focus()
	// // input.CursorStyle = lipgloss.Style{}
	// input.Cursor.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	// input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#1e90ff"))
	input.CharLimit = 156
	input.Width = 20
	
	return
}