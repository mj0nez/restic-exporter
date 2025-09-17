package collector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type CmdBuffer struct {
	Message string //`mapstructure:"message"`
}

func run(ctx context.Context, binPath string, args []string, repo string, password string, handleErrors bool) (*bytes.Buffer, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cwd, err := os.Getwd()
	if err != nil {
		return stdout, err
	}
	env := make(map[string]string)

	env["RESTIC_REPOSITORY"] = repo
	env["RESTIC_PASSWORD"] = password

	command := prepareCommand(ctx, binPath, cwd, args, env, stdout, stderr)

	// run, handle errors and parse the returned message
	err = command.Run()
	if handleErrors && err != nil {
		cmdErr := &CmdBuffer{}
		if err := json.Unmarshal(stderr.Bytes(), cmdErr); err != nil {
			fmt.Printf("Failed to parse command output! Raw message: %v", stderr.String())
		} else {
			fmt.Println(cmdErr.Message)
		}
	}

	return stdout, err
}

func prepareCommand(ctx context.Context, binPath string, cwd string, args []string, env map[string]string, stdout io.Writer, stderr io.Writer) *exec.Cmd {
	var cmd *exec.Cmd

	if ctx != nil {
		cmd = exec.CommandContext(ctx, binPath, args...)
	} else {
		cmd = exec.Command(binPath, args...)
	}
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

	return cmd
}

// TODO consider extracting the common parts of ech cli usage to util func
// func runResticCli() {}
