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
	cursor int
	input textinput.Model
	tables []table.Model
    // selected map[int]struct{}
	search bool
	err error
}

func initialModel() model {
	return model{
		cursor: 0,
		input: custom.SearchBar(),
		tables:  []table.Model{docker.TableContainers(), docker.TableImages()},
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
				m.tables[m.cursor].Focus()
			default:
				m.input, cmd = m.input.Update(msg)
			}
		} else {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "ctrl+a":
				m.cursor = 0
			case "ctrl+e":
				m.cursor = 1
			case ":":
				m.search = true
				m.tables[m.cursor].Blur()
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
	var heightToRemove int
	if (m.search) {
		heightToRemove = 2 + 2 + inputHeight + descHeight
	} else {
		heightToRemove = 2 + descHeight
	}
	widthToRemove := 2

	inputStyle := custom.CreateStyle(width-widthToRemove, inputHeight, "#ffa500")
	tableStyle := custom.CreateStyle(width-widthToRemove, height-heightToRemove, "12")

	// fmt.Println(width)
	// option := table.WithWidth(width)
	// m.table.SetWidth(width-widthToRemove)
	// (&m.table).UpdateViewport()
	// m.table.UpdateViewport()

	input := inputStyle.Render(m.input.View())
	table := tableStyle.Render(m.tables[m.cursor].View())
	desc := docker.Lists()

	var result string
	if (m.search) {
		result = lipgloss.JoinVertical(lipgloss.Left, desc, input, table)
	} else {
		result = lipgloss.JoinVertical(lipgloss.Left, desc, table)
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