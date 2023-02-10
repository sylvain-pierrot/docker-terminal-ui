package main

import (
	"context"
	"fmt"
	"log"
	"main/custom"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)



var rows_images []table.Row
var rows_containers []table.Row
var columns_images = []table.Column{
	{Title: "REPOSITORY", Width: 10},
	{Title: "TAG", Width: 10},
	{Title: "IMAGE ID", Width: 20},
	{Title: "SIZE", Width: 10},
}

var columns_containers = []table.Column{
	{Title: "CONTAINER ID", Width: 15},
	{Title: "IMAGE", Width: 10},
	{Title: "COMMAND", Width: 15},
	{Title: "CREATED", Width: 15},
	{Title: "STATUS", Width: 15},
	{Title: "PORTS", Width: 10},
	{Title: "NAMES", Width: 15},
}
	
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
			m.table = custom.CreateTable(rows_containers, columns_containers)
		case "ctrl+e":
			m.table = custom.CreateTable(rows_images, columns_images)
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

	inputStyle := createStyle(1, width-widthTotalMargin, "#ffa500")
	tableStyle := createStyle(height-inputStyle.GetHeight()-heightTotalMargin, width-widthTotalMargin, "12")

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

	// containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		port :=  strconv.FormatUint(uint64(container.Ports[0].PrivatePort), 10) + "/" + container.Ports[0].Type
		row := []string{container.ID[:12], container.Image, container.Command, string(rune(container.Created)), container.Status, port, strings.Trim(container.Names[0], "/")	}
		rows_containers = append(rows_containers, row)	
	}

	// images
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		repo_tag := strings.Split(image.RepoTags[0], ":")
		image_id := strings.Split(image.ID, ":")
		row := []string{repo_tag[0], repo_tag[1], image_id[1][:12], strconv.FormatInt(image.Size, 10)}
		rows_images = append(rows_images, row)
		// fmt.Printf("%s 	%s %s %dB\n", repository, tag, image.ID, image.Size)
	}

	// create table
	t := custom.CreateTable(rows_images, columns_images)

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

func createStyle(height, width int, color string) lipgloss.Style {
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