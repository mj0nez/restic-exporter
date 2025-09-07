package internal

import (
	"context"
	"time"
)

func StartWorker(ctx context.Context, interval time.Duration, binPath string, repos []string) error {

	wCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// prefetch
	for _, repo := range repos {
		go collect(wCtx, binPath, repo)
	}

	for {
		select {
		case <-wCtx.Done():
			// exiting
			return nil
		case <-time.After(interval):
			for _, repo := range repos {
				go collect(wCtx, binPath, repo)
			}
		}
	}
}
