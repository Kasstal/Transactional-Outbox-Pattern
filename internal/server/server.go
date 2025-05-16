package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"orders-center/internal/utils"
	"time"
)

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func NewServer(config utils.Config, router *gin.Engine) *Server {

	srv := &http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}

	return &Server{
		httpServer:      srv,
		shutdownTimeout: config.ServerShutDownTime,
	}
}

func (s *Server) Run() {

	go func() {
		fmt.Println("Gin server listening on :8080")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: %v\n", err)

		}
	}()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error: %v\n", err)
	}
}
