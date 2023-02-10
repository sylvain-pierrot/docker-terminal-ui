package main

import (
	"fmt"
	"log"
	"main/custom"
	"main/docker"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/client"
)

type model struct {
	table table.Model
	textinput textinput.Model
	err       error
}

// Init
func (m model) Init() tea.Cmd { 
	// return tea.Batch()
	return tea.Batch(tea.EnterAltScreen)

}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+a":
			m.table = custom.CreateTable(docker.ListContainers(cli))
		case "ctrl+e":
			m.table = custom.CreateTable(docker.ListImages(cli))
		case ":":
			m.textinput.Focus()
			switch msg.Type {
			case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			}		
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

// view
func (m model) View() string {
	height, width := getWindowSize()

	widthTotalMargin := 2
	heightTotalMargin := 4

	inputStyle := custom.CreateStyle(1, width-widthTotalMargin, "#ffa500")
	tableStyle := custom.CreateStyle(height-inputStyle.GetHeight()-heightTotalMargin, width-widthTotalMargin, "12")

	input := inputStyle.Render(m.textinput.View())
	table := tableStyle.Render(m.table.View())
	
	number := tableStyle.GetHeight()

	return lipgloss.PlaceVertical(height, lipgloss.Top, lipgloss.JoinVertical(lipgloss.Left, input, table) + "\n" + fmt.Sprint(height) + " " + strconv.Itoa(number) )
	
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// create table
	t := custom.CreateTable(docker.ListImages(cli))

	// create input
	ti := textinput.New()
	// ti.Placeholder = "Pikachu"
	// ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	// add table to model
	m := model{t, ti, nil}
	// launch
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func getWindowSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	output := string(out)
	dimensions := strings.Split(strings.TrimSpace(output), " ")
	height, err := strconv.Atoi(dimensions[0])
	width, err := strconv.Atoi(dimensions[1])

	if err != nil {
		log.Fatal(err)
	}

	return height, width
}

