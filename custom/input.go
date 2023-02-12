package custom

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func SearchBar() (input textinput.Model) {
	input = textinput.New()
	input.Prompt ="🐨>"
	// ti.Focus()
	input.CharLimit = 156
	input.Width = 20
	
	return
}