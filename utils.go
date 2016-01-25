package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunCommand(c string) {
	words := strings.Fields(c)
	name := words[0]
	args := words[1:]
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(output), time.Now().Second())
}
