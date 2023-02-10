package docker

import (
	"context"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ListContainers(cli *client.Client) (rows []table.Row, columns []table.Column) {
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
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		port :=  strconv.FormatUint(uint64(container.Ports[0].PrivatePort), 10) + "/" + container.Ports[0].Type
		row := []string{container.ID[:12], container.Image, container.Command, string(rune(container.Created)), container.Status, port, strings.Trim(container.Names[0], "/")	}
		rows = append(rows, row)	
	}

	return
}

func ListImages(cli *client.Client) (rows []table.Row, columns []table.Column) {
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

	return
}