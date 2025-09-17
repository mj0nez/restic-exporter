package collector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mj0nez/restic-exporter/contrib/restic"
	"github.com/mj0nez/restic-exporter/internal/config"
)

const statsMode string = "raw-data"

// var ErrCheck = fmt.Errorf("repository check failed")

// func boolToInt(val bool) int8 {
// 	var i int8
// 	if val {
// 		return 1
// 	}
// 	return i
// }

func GetStats(ctx context.Context, binPath string, repo config.Repository) {
	_, err := checkRepo(ctx, binPath, repo.Restic.Repo, repo.Restic.Password)

	if err != nil {
		if errors.Is(err, ErrCheck) {
			slog.Error(fmt.Sprintf("Failed to get stats in repo %v because: %v", repo, err))
			return
		}
	}
}

func getStats(ctx context.Context, binPath string, repo string, password string) (*restic.StatsContainer, error) {
	// check and verify integrity fo the repository
	stats := &restic.StatsContainer{}

	// TODO: consider adding an option to change the mode:
	// https://restic.readthedocs.io/en/stable/manual_rest.html#getting-information-about-repository-data
	args := []string{"--no-lock", "stats", "--mode", statsMode, "--json"}

	stdout, err := run(ctx, binPath, args, repo, password, true)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(stdout.Bytes(), stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
