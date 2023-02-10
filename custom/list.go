package custom

import (
	"github.com/charmbracelet/lipgloss"
)

func Lists() string {
	cmdColor := lipgloss.Color("#1e90ff")
	descColor := lipgloss.Color("#D9DCCF")
	infoColor := lipgloss.Color("#ffa500")

	listInfo := lipgloss.NewStyle().
		Foreground(infoColor).
		// MarginRight(1).
		Height(5).
		Width(10).
		Bold(true)

	listCmd := lipgloss.NewStyle().
		Foreground(cmdColor).
		// MarginRight(1).	
		Height(5).
		Width(10).
		Bold(true)
	
	listDesc := lipgloss.NewStyle().
		Foreground(descColor).
		MarginRight(10).
		Height(5).
		Width(10).
		Bold(true)

		
		// 383838
	listItem := lipgloss.NewStyle().PaddingLeft(2).Render

	lists2 := lipgloss.JoinHorizontal(lipgloss.Right,
		listCmd.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listItem("<ctrl-d>"),
				listItem("<d>"),
				listItem("<e>"),
				listItem("<?>"),
				listItem("<u>"),
				listItem("<y>"),
			),
		),
		listDesc.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listItem("Delete"),	
				listItem("Describe"),
				listItem("Edit"),
				listItem("Help"),
				listItem("Use"),
				listItem("YAML"),
			),
		),
	)

	lists1 := lipgloss.JoinHorizontal(lipgloss.Right,
		listInfo.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listItem("Context"),
				listItem("Cluster"),
				listItem("K9s Rev"),
				listItem("K8s Rev"),
				listItem("CPU"),
				listItem("MEM"),
			),
		),
		listDesc.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listItem("Delete"),	
				listItem("Describe"),
				listItem("Edit"),
				listItem("Help"),
				listItem("Use"),
				listItem("YAML"),
			),
		),
	)

	// doc := strings.Builder{}
	// doc.WriteString(lists1 + "\n\n")
	// doc.WriteString(lists2)

	docStyle := lipgloss.NewStyle().Padding(1, 2, 1, 2)

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, lists1, lists2)) 
	// return docStyle.Render(doc.String())
	// return doc.String()
}