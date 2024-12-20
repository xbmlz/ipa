package ipa

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/xbmlz/ipa/config"
	"github.com/xbmlz/ipa/logger"
)

type App struct {
	Config config.Config

	httpServer *httpServer

	httpRegistered bool
}

func New() *App {
	app := &App{}

	// Load config
	app.loadConfig()

	// HTTP server
	host := app.Config.GetString("HTTP_HOST", defaultHTTPHost)
	port := app.Config.GetInt("HTTP_PORT", defaultHTTPPort)
	app.httpServer = newHTTPServer(host, port)

	return app
}

// loadConfig loads the configuration from the default path.
func (a *App) loadConfig() {
	var configPath string
	if _, err := os.Stat(config.DefaultConfigPath); err == nil {
		configPath = config.DefaultConfigPath
	}
	a.Config = config.New(configPath, logger.New("INFO"))
}

// Run starts the application.
func (a *App) Run() {
	// Create a context that is canceled on receiving termination signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Goroutine to handle shutdown when context is canceled
	go func() {
		<-ctx.Done()

		// Create a shutdown context with a timeout
		shutdownCtx, done := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer done()

		_ = a.Shutdown(shutdownCtx)
	}()

	wg := sync.WaitGroup{}

	// Start HTTP Server
	if a.httpRegistered {
		wg.Add(1)

		go func(s *httpServer) {
			defer wg.Done()
			s.Run()
		}(a.httpServer)
	}

	wg.Wait()
}

// Shutdown gracefully shuts down the application.
func (a *App) Shutdown(ctx context.Context) error {
	var err error
	if a.httpServer != nil {
		err = errors.Join(err, a.httpServer.Shutdown(ctx))
	}
	return err
}

func (a *App) addRoute(method, path string, handler HandlerFunc) {
	a.httpRegistered = true
	a.httpServer.router.Handle(method, path, func(c *gin.Context) {
		ctx := &Context{
			Context: c,
			Log:     logger.New("INFO"),
		}
		handler(ctx)
	})
}

// Get
func (a *App) Get(path string, handler HandlerFunc) {
	a.addRoute("GET", path, handler)
}
