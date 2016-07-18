package ps4control

import (
	"os/exec"
	"strings"
	"fmt"
)

const (
	ps4tool = "ps4-wake"
	args = "-vPB"
)

func runCommand() (string,  error) {
	command := exec.Command(ps4tool, args)
	text, err := command.CombinedOutput()
	return string(text), err
}

const (
	StatePlaying = iota
	StatePaused
	StateStopped
)

func GetPs4State() int {
	output, _ := runCommand()
	fmt.Printf("ps4: %v\n", output)
	if output == "" {
		return StateStopped
	}

	if strings.Contains(output, "Device found") {
		fmt.Printf("have device found")
		if strings.Contains(output, "Standby") {
			return StateStopped
		}

		if strings.Contains(output, "Home") {
			return StatePaused
		}

		return StatePlaying
	}

	return StateStopped
}



