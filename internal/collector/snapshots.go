package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mj0nez/restic-exporter/contrib/restic"
	"github.com/mj0nez/restic-exporter/internal/config"
	"github.com/mj0nez/restic-exporter/internal/metrics"
)

func Collect(binPath string, repo config.Repository) {
	snaps, err := getSnapshots(context.TODO(), binPath, repo.Restic.Repo, repo.Restic.Password)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %v", repo.Restic.Repo, err))
		return
	}

	fmt.Printf("%+v", snaps)

}

func GetSnapshots(ctx context.Context, binPath string, repo config.Repository) {
	snaps, err := getSnapshots(ctx, binPath, repo.Restic.Repo, repo.Restic.Password)

	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %+v", repo, err))
	}

	// metrics.CheckSuccess.
	metrics.SnapshotsTotal.WithLabelValues(repo.Name).Set(float64(len(snaps)))

}

func getSnapshots(ctx context.Context, binPath string, repo string, password string) ([]restic.Snapshot, error) {

	args := []string{"--no-lock", "snapshots", "--json"}

	stdout, err := run(ctx, binPath, args, repo, password, true)
	if err != nil {
		return nil, err
	}

	snapshots := make([]restic.Snapshot, 0, 10)

	err = json.Unmarshal(stdout.Bytes(), &snapshots)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}
