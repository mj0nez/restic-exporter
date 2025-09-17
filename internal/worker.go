package internal

import (
	"context"
	"time"

	"github.com/mj0nez/restic-exporter/internal/config"
)

type collectorFn func(ctx context.Context, binPath string, repo config.Repository)

func startWorker(collector collectorFn, ctx context.Context, prefetch bool, interval time.Duration, binPath string, repo config.Repository) {

	wCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if prefetch {
		go collector(wCtx, binPath, repo)
	}

	for {
		select {
		case <-wCtx.Done():
			return // exiting
		case <-time.After(interval):
			go collector(wCtx, binPath, repo)
		}
	}
}
