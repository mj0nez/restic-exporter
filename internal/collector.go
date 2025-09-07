package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/mj0nez/restic-exporter/contrib/restic"
	"github.com/mj0nez/restic-exporter/internal/metrics"
)

func Collect() {
	snaps, err := getSnapshots(nil, "restic", ".tmp/repo")
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %v", ".tmp/repo", err))
	}

	fmt.Printf("%v", snaps)

}

func collect(ctx context.Context, binPath string, repo string) {
	snaps, err := getSnapshots(ctx, "restic", repo)

	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %v", repo, err))
	}
	// this was a counter, but I assume with pruning the value might change
	// furthermore this avoids incrementing until we reach the curre
	metrics.SnapshotsTotal.WithLabelValues(repo).Set(float64(len(snaps)))

}

func getSnapshots(ctx context.Context, binPath string, repo string) ([]restic.Snapshot, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// create defaults
	args := []string{"-r", repo, "--no-lock", "snapshots", "--json"}

	env := make(map[string]string)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err = runCommand(ctx, binPath, cwd, args, env, stdout, stderr)
	if err != nil {
		return nil, err
	}
	// fmt.Println(stdout.String())
	// fmt.Println(stderr.String())

	snapshots := make([]restic.Snapshot, 0, 10)

	err = json.Unmarshal(stdout.Bytes(), &snapshots)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}
