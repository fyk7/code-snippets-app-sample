package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fyk7/code-snippets-app/app/config"
	"github.com/fyk7/code-snippets-app/app/di"
	_handler "github.com/fyk7/code-snippets-app/app/interface_adapter/handler"
	_middleware "github.com/fyk7/code-snippets-app/app/interface_adapter/handler/middleware"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.LoadConf()

	// Dependency Injection
	serviceContainer := di.Initialize(cfg, cfg.AppTimeOut)

	e := echo.New()
	mw := _middleware.InitMiddleware()
	e.Use(mw.CORS)

	// Register handlers.
	_handler.NewSnippetHandler(e, serviceContainer.SnippetService)
	_handler.NewTagHandler(e, serviceContainer.TagService)

	// Graceful shutdown with signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	// Start HTTP server
	g.Go(func() error {
		slog.Info("starting server", "addr", ":8080")
		return e.Start(":8080")
	})

	// Wait for shutdown signal, then drain gracefully in parallel
	g.Go(func() error {
		<-ctx.Done()
		slog.Info("shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.AppTimeOut)
		defer cancel()

		// Shut down HTTP server and close DB connection in parallel
		sg, _ := errgroup.WithContext(shutdownCtx)

		sg.Go(func() error {
			return e.Shutdown(shutdownCtx)
		})

		sg.Go(func() error {
			sqlDB, err := serviceContainer.DB.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		})

		return sg.Wait()
	})

	if err := g.Wait(); err != nil {
		slog.Info("server stopped", "reason", err.Error())
	}
}
