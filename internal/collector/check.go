package collector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/mj0nez/restic-exporter/contrib/restic"
	"github.com/mj0nez/restic-exporter/internal/config"
	"github.com/mj0nez/restic-exporter/internal/metrics"
)

var ErrCheck = fmt.Errorf("repository check failed")

func boolToInt(val bool) int8 {
	var i int8
	if val {
		return 1
	}
	return i
}

func RunCheck(ctx context.Context, binPath string, repo config.Repository) {
	check, err := checkRepo(ctx, "restic", repo.Restic.Repo, repo.Restic.Password)

	if err != nil {
		if errors.Is(err, ErrCheck) {
			slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %v", repo, err))
		} else {
			metrics.CheckFailed.WithLabelValues(repo.Name).Inc()
		}
	} else {
		// metrics.CheckSuccess.
		metrics.CheckSuccess.WithLabelValues(repo.Name).Inc()
	}

	metrics.CheckSuggestRepairIndex.WithLabelValues(repo.Name).Set(float64(boolToInt(check.HintRepairIndex)))
	metrics.CheckSuggestPrune.WithLabelValues(repo.Name).Set(float64(boolToInt(check.HintPrune)))
	metrics.CheckErrorsTotal.WithLabelValues(repo.Name).Set(float64(check.NumErrors))

}

func checkRepo(ctx context.Context, binPath string, repo string, password string) (*restic.CheckSummary, error) {
	// check and verify integrity fo the repository
	summary := &restic.CheckSummary{}

	args := []string{"--no-lock", "check", "--json"}

	// this command uses the return code to indicate the repository's status
	// therefore rc != 0 should not be logged
	stdout, err := run(ctx, binPath, args, repo, password, false)

	// TODO: check if it would be better to add a rc allow-list to the run function
	// 		I'm currently unsure if that would be wise, because there would be no
	//		clear way to map logic to each code - like below...
	if err != nil {
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
