package zfs

import (
	"context"
	"os/exec"
	"time"
)

var execTimeout = time.Second * 60
var execShell = "/usr/bin/bash"

func executeScript(script []byte, args ...string) ([]byte, error) {

	// setup context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	// execute the bash script
	// the -c interprets the input string as script
	cmd := exec.CommandContext(ctx, execShell, append([]string{"-c", string(script)}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		// errored
		return nil, err
	}

	// no error, success
	return out, nil
}
