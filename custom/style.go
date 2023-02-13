package custom

import (
	"github.com/charmbracelet/lipgloss"
)	 

func CreateStyle(width, height int, color string) lipgloss.Style {
	style := lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color(color)).
	Foreground(lipgloss.Color("#7accf8")).
    PaddingLeft(1).
    PaddingRight(1).
	Width(width)

	if (height > 0) {
		style.Height(height)
	}

	return style
}