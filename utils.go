package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunCommand(c string) {
	words := strings.Fields(c)
	cmd, _ := CreateCommand(words)
	output, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))
}

func CreateCommand(args []string) (*exec.Cmd, error) {
	name := args[0]
	cmd := &exec.Cmd{
		Path: name,
		Args: args,
	}
	if filepath.Base(name) == name {
		lp, err := exec.LookPath(name)
		if err != nil {
			return cmd, err
		}

		cmd.Path = lp
	}
	return cmd, nil
}
