package internal

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func runCommand(binPath string, cwd string, args []string, env map[string]string, stdout io.Writer, stderr io.Writer) error {

	cmd := exec.Command(binPath, args...)

	// handle environment and overrides
	// we create a working copy, replace the overrides and rebuild the required array
	os_env := os.Environ()
	cmd_env := make(map[string]string, len(os_env))

	for _, envVar := range os_env {
		keyVal := strings.Split(envVar, "=")
		cmd_env[keyVal[0]] = keyVal[1]
	}
	// set overrides
	for k, v := range env {
		cmd_env[k] = v
	}
	// rebuild
	cmd_env_res := make([]string, len(cmd_env))
	for k, v := range cmd_env {
		cmd_env_res = append(cmd_env_res, k+"="+v)
	}
	cmd.Env = cmd_env_res

	cmd.Dir = cwd
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}
