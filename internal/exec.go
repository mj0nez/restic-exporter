package internal

import (
	"io"
	"os"
	"os/exec"
)

func runCommand(binPath string, cwd string, args []string, env map[string]string, stdout io.Writer, stderr io.Writer) error {

	cmd := exec.Command(binPath, args...)
	cmd.Env = os.Environ()

	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Dir = cwd
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}
