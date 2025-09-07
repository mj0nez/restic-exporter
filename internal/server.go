package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var signalsToListenTo = []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM}

type HttpServerOpts struct {
	TLSConfig         *tls.Config
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func NewHttpServerOpts() HttpServerOpts {
	return HttpServerOpts{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func NewHttpServer(addr string, handler http.Handler, opts HttpServerOpts) *http.Server {

	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		TLSConfig:         opts.TLSConfig,
		ReadTimeout:       opts.ReadTimeout,
		ReadHeaderTimeout: opts.ReadHeaderTimeout,
		WriteTimeout:      opts.WriteTimeout,
		IdleTimeout:       opts.IdleTimeout,
	}
}

func NewRouter(registry *prometheus.Registry) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Expose /metrics HTTP endpoint using the created custom registry.
	router.GET("/metrics", func(ctx *gin.Context) {
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}).ServeHTTP(ctx.Writer, ctx.Request)
	})

	return router
}

func RunServer(server *http.Server) error {
	// creates a new goroutine which calls context.cancel as soon
	// as any of signalsToListenTo is received
	ctx, stop := signal.NotifyContext(context.Background(), signalsToListenTo...)

	// ensure stop is called when Run returns
	// NotifyContext starts background routine for the signal handling
	defer stop()

	log := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	log.Debug("Starting background worker.")

	go startWorker(collectSnapshots, ctx, time.Duration(15)*time.Second, "restic", []string{".tmp/repo", ".tmp/repo2"})
	go startWorker(collectCheck, ctx, time.Duration(20)*time.Second, "restic", []string{".tmp/repo", ".tmp/repo2"})

	log.Debug("Starting http listener.")
	log.Info(fmt.Sprintf("Listening on %s", server.Addr))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(fmt.Sprintf("listen: %v", err))
		}
	}()

	<-ctx.Done()

	log.Info("Received Shutdown signal. Stopping agent gracefully. Press Ctrl+C to force.")
	stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("Agent forced to shutdown: %v", err))
	}
	log.Info("Agent exiting")

	return nil
}
