package custom

import "github.com/charmbracelet/lipgloss"	 

func CreateStyle(height, width int, color string) lipgloss.Style {
	style := lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
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