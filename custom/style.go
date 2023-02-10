package custom

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)	 

func CreateStyle(width, height int, color string, table table.Model) lipgloss.Style {
	var sumWidthCells int
	for i := 0; i < len(columns); i++ {
		sumWidthCells += columns[i].Width
	}
	if sumWidthCells > windowWidth

	style := lipgloss.NewStyle().
	// BorderStyle(lipgloss.NormalBorder()).
	// BorderForeground(lipgloss.Color(color)).
	Foreground(lipgloss.Color("#7accf8")).
	// PaddingTop(1).
    // PaddingBottom(1).
    PaddingLeft(1).
    PaddingRight(1)
	// Width(width)

	if (height > 0) {
		style.Height(height)
	}

	return style
}