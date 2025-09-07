package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mj0nez/restic-exporter/contrib/restic"
)

func Collect() ([]restic.Snapshot, error) {
	snaps, err := collect(".tmp/repo")
	return snaps, err
}

func collect(repo string) ([]restic.Snapshot, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// create defaults
	args := []string{"--no-lock", "snapshots", "--json"}

	env := make(map[string]string)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err = runCommand("restic", cwd, args, env, stdout, stderr)
	if err != nil {
		return nil, err
	}
	fmt.Println(stdout.String())
	fmt.Println(stderr.String())

	snapshots := make([]restic.Snapshot, 0, 10)

	err = json.Unmarshal(stdout.Bytes(), &snapshots)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

func CollectAll() error {
	return nil
}
