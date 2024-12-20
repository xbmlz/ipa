package ipa

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultHTTPHost = "localhost"
	defaultHTTPPort = 8000
)

type httpServer struct {
	router *gin.Engine
	host   string
	port   int
	srv    *http.Server
}

func newHTTPServer(host string, port int) *httpServer {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	return &httpServer{
		router: r,
		host:   host,
		port:   port,
	}
}

func (s *httpServer) Run() {
	if s.srv != nil {
		slog.Warn("Server already running on port", slog.Int("port", s.port))
		return
	}

	slog.Info("Starting HTTP server on %s:%d", s.host, s.port)

	s.srv = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.host, s.port),
		Handler:           s.router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// If no certFile/keyFile is provided, run the HTTP server
	if err := s.srv.ListenAndServe(); err != nil {
		slog.Error("Failed to start HTTP server", slog.String("error", err.Error()))
	}
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	return ShutdownWithContext(ctx, func(ctx context.Context) error {
		return s.srv.Shutdown(ctx)
	}, func() error {
		if err := s.srv.Close(); err != nil {
			return err
		}

		return nil
	})
}
