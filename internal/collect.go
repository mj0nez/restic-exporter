package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mj0nez/restic-exporter/contrib/restic"
)

func Collect() error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// create defaults
	args := []string{"snapshots", "--json"}
	env := make(map[string]string)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err = runCommand("restic", cwd, args, env, stdout, stderr)
	if err != nil {
		return err
	}
	// fmt.Println(stdout.String())
	// fmt.Println(stderr.String())

	snapshots := make([]restic.Snapshot, 0, 10)

	err = json.Unmarshal(stdout.Bytes(), &snapshots)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", snapshots[0].Summary)

	return nil
}
