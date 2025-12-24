// Package shell provides utilities for executing shell commands.
package shell

import (
	"fmt"
	"os/exec"
)

// ExecuteShellCommand runs a command in a bash shell and returns any error.
func ExecuteShellCommand(command string) error {
	// Create the command and pass it to a new shell
	cmd := exec.Command("bash", "-c", command)

	// Run the command and get any output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %v\nOutput: %s", err, string(output))
	}

	// Print the output if the command was successful
	// fmt.Println("Output:", string(output))
	return nil
}
