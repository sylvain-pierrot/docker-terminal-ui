package custom

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func CreateTable(rows []table.Row, columns []table.Column) table.Model {
	width := 0
	for i := 0; i < len(columns); i++ {
		width += columns[i].Width
	}
	
	t := table.New(
		table.WithRows(rows),
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(5),
		table.WithWidth(width),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#87cefa")).
		Bold(true)
	t.SetStyles(s)
	
	return t
}