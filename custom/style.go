package custom

import (
	"github.com/charmbracelet/lipgloss"
)	 

func CreateStyle(width, height int, color string) lipgloss.Style {
	// var sumWidthCells int
	// for i := 0; i < len(columns); i++ {
	// 	sumWidthCells += columns[i].Width
	// }
	// if sumWidthCells > windowWidth

	style := lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color(color)).
	Foreground(lipgloss.Color("#7accf8")).
	// PaddingTop(1).
    // PaddingBottom(1).
    PaddingLeft(1).
    PaddingRight(1).
	Width(width)

	if (height > 0) {
		style.Height(height)
	}

	return style
}