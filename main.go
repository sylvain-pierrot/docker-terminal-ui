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
    selected map[int]struct{}
	inputSelected bool
	// table table.Model
	err       error
}

func initialModel() model {
	return model{
		cursor: 1,
		input: custom.SearchBar(),
		tables:  []table.Model{docker.TableContainers, docker.TableImages},
		selected: make(map[int]struct{}),
		inputSelected: false,
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
		switch msg.String() {
		case "esc":
			m.inputSelected = false
			m.textinput.Blur()
			m.table.Focus()
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+a":
			m.inputSelected = false
			m.table = custom.CreateTable(docker.ListContainers())
		case "ctrl+e":
			m.inputSelected = false
			m.table = custom.CreateTable(docker.ListImages())
		case ":":
			m.inputSelected = true
			m.table.Blur()
			m.textinput.Focus()
			// switch msg.Type {
			// case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			// 	return m, tea.Quit
			// }		
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.tables[i], cmd = m.tables[i].Update(msg)
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// view
func (m model) View() string {
	// legend := "<ctrl-a> Containers\n<ctrl-e> Images\n"

	// lipgloss.lipgloss.NewStyle()
	height, width := utils.GetWindowSize()
	inputHeight := 1
	var heightToRemove int
	if (m.inputSelected) {
		heightToRemove = 2 + 2 + inputHeight
	} else {
		heightToRemove = 2
	}
	widthToRemove := 2



	inputStyle := custom.CreateStyle(width-widthToRemove, inputHeight, "#ffa500")
	tableStyle := custom.CreateStyle(width-widthToRemove, height-heightToRemove, "12")

	// fmt.Println(width)
	// option := table.WithWidth(width)
	// m.table.SetWidth(width-widthToRemove)
	// (&m.table).UpdateViewport()
	// m.table.UpdateViewport()

	input := inputStyle.Render(m.textinput.View())
	table := tableStyle.Render(m.table.View())


	var result string
	if (m.inputSelected) {
		result = lipgloss.JoinVertical(lipgloss.Top, input, table)
	} else {
		result = table
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



