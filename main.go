package main

import (
	"fmt"
	"main/custom"
	"main/docker"
	"main/utils"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/client"
)

type model struct {
	cursor string
	input textinput.Model
    tables map[string]table.Model
	tableAction table.Model
	action bool
	search bool
	err error
}

func initialModel() model {
	return model{
		cursor: "container",
		input: custom.SearchBar(),
		tables: map[string]table.Model{"container": docker.TableContainers(), "image": docker.TableImages(), "volume": docker.TableVolumes(), "network": docker.TableNetworks()},
		tableAction: table.Model{},
		action: false,
		search: false,
		err: nil,
	}
}

// Init
func (m model) Init() tea.Cmd { 
	return tea.Batch(tea.EnterAltScreen)
}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	options := []string{"container", "image", "volume", "network"}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if m.search {
			switch msg.String() {

			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.input.Reset()
				m.search = false
				m.input.Blur()
				// m.tables[m.cursor].Focus()
				// table.Focus()
			case "enter":
				if _, exists := m.tables[m.input.Value()]; exists {
					m.cursor = m.input.Value()
				}
			case "tab":
				m.input.CursorEnd()
			default:
				m.input, cmd = m.input.Update(msg)

				pos := m.input.Position()
				val := m.input.Value()[:pos]

				m.input.SetValue(val)
				m.input.SetCursor(len(val))

				var completions []string
				for _, option := range options {
					if strings.HasPrefix(option, val) {
						completions = append(completions, option)
					}
				}

				if len(completions) > 0 && len(m.input.Value()) > 0 {
					m.input.Reset()
					m.input.SetValue(completions[0])
					m.input.SetCursor(pos)
				}
			}
		} else if m.action {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.tableAction = table.Model{}
				m.action = false
				m.tableAction.Blur()
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "i":
				row := m.tables[m.cursor].SelectedRow()
				containerID := row[0]
				m.tableAction = docker.TableContainerInspect(containerID)
				m.action = true
				m.tableAction.Focus()
			case ":":
				m.search = true
				m.input.Focus()
			case "enter":
				return m, tea.Batch(
					tea.Printf("Let's go to !"),
				)
			}
		}
	}
	m.tables[m.cursor], cmd = m.tables[m.cursor].Update(msg)
	

	return m, cmd
}

// view
func (m model) View() string {
	// width & height
	height, width := utils.GetWindowSize()

	contextHeight := 1
	descHeight := 16
	inputHeight := 1
	var tableHeight int

	if (m.search) {
		tableHeight = height - inputHeight - descHeight - contextHeight - 2
	} else {
		tableHeight = height - descHeight - contextHeight
	}
	width -= 2

	if (tableHeight < 1) {
		tableHeight = 0
	}
	
	table := m.tables[m.cursor]
	currentWidth := table.Width() + 18

	if (currentWidth <= width) {
		currentWidth = width
	}

	inputStyle := custom.CreateStyle(width, inputHeight, "#ffa500")
	tableStyle := custom.CreateStyle(currentWidth, tableHeight, "12")

	table.SetWidth(currentWidth)
	table.SetHeight(tableHeight)

	input_string := inputStyle.Render(m.input.View())
	table_string := tableStyle.Render(table.View())
	desc_string := docker.Lists()
	logo_string := docker.LogoDO3()
	legend_string := lipgloss.JoinHorizontal(lipgloss.Left, desc_string, logo_string)
	context_string := docker.LabelContext(m.cursor)

	var result string
	if (m.search) {
		result = lipgloss.JoinVertical(lipgloss.Left, legend_string, input_string, table_string, context_string)
	} else if m.action {
		m.tableAction.SetWidth(currentWidth)
		m.tableAction.SetHeight(tableHeight)
		tableAction_string := tableStyle.Render(m.tableAction.View())
		result = lipgloss.JoinVertical(lipgloss.Left, legend_string,  tableAction_string, context_string)
	} else {		
		result = lipgloss.JoinVertical(lipgloss.Left, legend_string,  table_string, context_string)
	}

	return lipgloss.PlaceVertical(height, lipgloss.Top, result)
}

func main() {
	// client docker
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	docker.Init(cli)

	// launch
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}