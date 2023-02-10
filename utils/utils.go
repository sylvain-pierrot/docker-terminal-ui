package utils

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetWindowSize() (int, int) {
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