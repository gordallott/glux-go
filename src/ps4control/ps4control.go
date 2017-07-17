package ps4control

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	ps4tool = "ps4-wake"
	args    = "-vPB"
)

func runCommand() (string, error) {
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
		v := strings.ToLower(output)
		switch {
		case strings.Contains(v, "standby"):
			return StateStopped
		case strings.Contains(v, "home"):
			fallthrough
		case strings.Contains(v, "youtube"):
			return StatePaused
		default:
			return StatePlaying
		}
	}

	return StateStopped
}
