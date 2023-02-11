package main

import (
	"fmt"
	"main/custom"
	"main/docker"
	"main/utils"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/client"
)

type model struct {
	cursor string
	input textinput.Model
	// tables []table.Model
    tables map[string]table.Model
	search bool
	err error
}

func initialModel() model {
	return model{
		cursor: "container",
		input: custom.SearchBar(),
		// tables:  []table.Model{docker.TableContainers(), docker.TableImages()},
		tables: map[string]table.Model{"container": docker.TableContainers(), "image": docker.TableImages()},
		// selected: make(map[int]struct{}),
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
			default:
				m.input, cmd = m.input.Update(msg)
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "ctrl+a":
				m.cursor = "container"
			case "ctrl+e":
				m.cursor = "image"
			case ":":
				m.search = true
				// m.tables[m.cursor].Blur()
				m.input.Focus()
				// switch msg.Type {
				// case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
				// 	return m, tea.Quit
				// }		
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
	// legend := "<ctrl-a> Containers\n<ctrl-e> Images\n"
	// lipgloss.lipgloss.NewStyle()

	height, width := utils.GetWindowSize()

	descHeight := 15
	inputHeight := 1
	var tableHeight int

	if (m.search) {
		tableHeight = height - inputHeight - descHeight - 4
	} else {
		tableHeight = height - descHeight - 2
	}
	width -= 2

	if (tableHeight < 1) {
		tableHeight = 0
	}
	inputStyle := custom.CreateStyle(width, inputHeight, "#ffa500")
	tableStyle := custom.CreateStyle(width, tableHeight, "12")

	table := m.tables[m.cursor]
	table.SetHeight(tableHeight)
	table.SetWidth(90)
	// table.SetWidth(width-widthToRemove)

	input_string := inputStyle.Render(m.input.View())
	table_string := tableStyle.Render(table.View())
	desc_string := docker.Lists()

	var result string
	if (m.search) {
		result = lipgloss.JoinVertical(lipgloss.Left, desc_string, input_string, table_string)
	} else {
		result = lipgloss.JoinVertical(lipgloss.Left, desc_string, table_string)
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