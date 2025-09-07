package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/mj0nez/restic-exporter/contrib/restic"
	"github.com/mj0nez/restic-exporter/internal/metrics"
)

var ErrCheck = fmt.Errorf("repository check failed")

func Collect() {
	snaps, err := getSnapshots(context.TODO(), "restic", ".tmp/repo")
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %v", ".tmp/repo", err))
	}

	fmt.Printf("%v", snaps)

}

func collectSnapshots(ctx context.Context, binPath string, repo string) {
	snaps, err := getSnapshots(ctx, "restic", repo)

	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %v", repo, err))
	}

	// metrics.CheckSuccess.
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

func boolToInt(val bool) int8 {
	var i int8
	if val {
		return 1
	}
	return i
}

func collectCheck(ctx context.Context, binPath string, repo string) {
	check, err := checkRepo(ctx, "restic", repo)

	if err != nil {
		if errors.Is(err, ErrCheck) {
			slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %v", repo, err))
		} else {
			metrics.CheckFailed.WithLabelValues(repo).Inc()
		}
	} else {
		// metrics.CheckSuccess.
		metrics.CheckSuccess.WithLabelValues(repo).Inc()
	}

	metrics.CheckSuggestRepairIndex.WithLabelValues(repo).Set(float64(boolToInt(check.HintRepairIndex)))
	metrics.CheckSuggestPrune.WithLabelValues(repo).Set(float64(boolToInt(check.HintPrune)))
	metrics.CheckErrorsTotal.WithLabelValues(repo).Set(float64(check.NumErrors))

}

func checkRepo(ctx context.Context, binPath string, repo string) (*restic.CheckSummary, error) {
	// check and verify integrity fo the repository
	summary := &restic.CheckSummary{}

	cwd, err := os.Getwd()
	if err != nil {
		return summary, err
	}

	// create defaults
	args := []string{"-r", repo, "--no-lock", "check", "--json"}

	env := make(map[string]string)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err = runCommand(ctx, binPath, cwd, args, env, stdout, stderr)

	if err == nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// cmd.Run returned with an non-zero exit code
			if exitError.ExitCode() == 1 {
				// the integrity check failed and there are probably errors
				return summary, ErrCheck
			}
		}
		return summary, err // something else happened
	}

	err = json.Unmarshal(stdout.Bytes(), summary)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func runResticCli() {}
