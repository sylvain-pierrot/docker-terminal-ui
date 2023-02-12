package docker

import (
	"context"
	"main/custom"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cli *client.Client

func Init(client *client.Client) {
	cli = client
}

func TableContainers() table.Model {
	var rows []table.Row
	var columns []table.Column

	columns = []table.Column{
		{Title: "CONTAINER ID", Width: 15},
		{Title: "IMAGE", Width: 10},
		{Title: "COMMAND", Width: 15},
		{Title: "CREATED", Width: 15},
		{Title: "STATUS", Width: 15},
		{Title: "PORTS", Width: 10},
		{Title: "NAMES", Width: 15},
	}

	// containers
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		var port string
		if len(container.Ports) > 0 {
			port = strconv.FormatUint(uint64(container.Ports[0].PrivatePort), 10) + "/" + container.Ports[0].Type
		} else {
			port = "-"
		}
		row := []string{container.ID[:12], container.Image, container.Command, string(rune(container.Created)), container.Status, port, strings.Trim(container.Names[0], "/")}
		rows = append(rows, row)	
	}

	return custom.CreateTable(rows, columns)
}

func TableImages() table.Model {
	var rows []table.Row
	var columns []table.Column

	columns = []table.Column{
	{Title: "REPOSITORY", Width: 10},
	{Title: "TAG", Width: 10},
	{Title: "IMAGE ID", Width: 20},
	{Title: "SIZE", Width: 10},
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
		rows = append(rows, row)
	}

	return custom.CreateTable(rows, columns)
}

func Lists() string {
	info, err := cli.Info(context.Background())
	if err != nil {
		panic(err)
	}

	cmdColor := lipgloss.Color("#1e90ff")
	descColor := lipgloss.Color("#D9DCCF")
	infoColor := lipgloss.Color("#ffa500")

	listInfo := lipgloss.NewStyle().
		Foreground(infoColor).
		// MarginRight(1).
		Height(10).
		Width(20).
		Bold(true)

	listCmd := lipgloss.NewStyle().
		Foreground(cmdColor).
		// MarginRight(1).	
		Height(10).
		Width(10).
		Bold(true)
	
	listDesc := lipgloss.NewStyle().
		Foreground(descColor).
		MarginRight(12).
		Height(10).
		Width(20).
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
				listItem("Server Version:"),
				listItem("Containers:"),
				listItem("Images:"),
				listItem("Runtimes:"),
				listItem("Kernel Version:"),
				listItem("Operating System:"),
				listItem("OSType:"),
				listItem("Architecture:"),
				listItem("CPUs:"),
			),
		),
		listDesc.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				listItem(info.ServerVersion),
				listItem(strconv.Itoa(info.Containers)),
				listItem(strconv.Itoa(info.Images)),
				listItem(info.DefaultRuntime),
				listItem(info.KernelVersion),
				listItem(info.OperatingSystem),
				listItem(info.OSType),
				listItem(info.Architecture),
				listItem(strconv.Itoa(info.NCPU)),
			),
		),
	)

	// doc := strings.Builder{}
	// doc.WriteString(lists1 + "\n\n")
	// doc.WriteString(lists2)

	docStyle := lipgloss.NewStyle().Padding(1, 0, 0, 0)

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, lists1, lists2)) 
	// return docStyle.Render(doc.String())
	// return doc.String()
}