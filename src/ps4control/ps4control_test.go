package ps4control

import (
	"testing"
	"fmt"
)

func TestRunCommand(t *testing.T) {
	stdout, err := runCommand()
	fmt.Printf("%v: %v\n", err, stdout)
}
