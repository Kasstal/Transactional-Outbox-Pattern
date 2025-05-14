package graceful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ShutdownConfig — конфиг для graceful shutdown.
type ShutdownConfig struct {
	Timeout  time.Duration
	Handlers []func() error
}

func NewShutDownConfig(timeout time.Duration, handlers ...func() error) *ShutdownConfig {
	return &ShutdownConfig{timeout, handlers}
}

func WaitForShutdown(server *http.Server, cfg ShutdownConfig) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down...\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v\n", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}

	for _, handler := range cfg.Handlers {
		err := handler()
		if err != nil {
			log.Printf("HTTP server shutdown error when running handlers: %v\n", err)
		}
	}

	select {
	case <-ctx.Done():
		log.Println("Graceful shutdown timeout exceeded. Forcing exit.")
	}
}
