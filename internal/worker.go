package internal

import (
	"context"
	"time"
)

type collector func(ctx context.Context, binPath string, repos string)

func startWorker(collector collector, ctx context.Context, interval time.Duration, binPath string, repos []string) {

	wCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// prefetch
	for _, repo := range repos {
		go collector(wCtx, binPath, repo)
	}

	for {
		select {
		case <-wCtx.Done():
			// exiting
			return
		case <-time.After(interval):
			for _, repo := range repos {
				go collector(wCtx, binPath, repo)
			}
		}
	}
}
