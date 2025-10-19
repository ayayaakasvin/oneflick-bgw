package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ayayaakasvin/trends-updater/internal/app"
	"github.com/ayayaakasvin/trends-updater/internal/config"
	"github.com/ayayaakasvin/trends-updater/internal/logger"
	"github.com/ayayaakasvin/trends-updater/internal/models/inner"

	// "github.com/ayayaakasvin/trends-updater/internal/repo/postgresql"
	// "github.com/ayayaakasvin/trends-updater/internal/repo/valkey"
	// "github.com/ayayaakasvin/trends-updater/internal/worker"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    cfg := config.MustLoadConfig()
    logger := logger.SetupLogger()
	_ = cfg
    shutdownChan := inner.NewShutdownChannel()
	go func() {
		msg := shutdownChan.Value()
		logger.Errorf("Error aquired: %v", msg)
		cancel()
	}()

    sigChannel := make(chan os.Signal, 1)
    signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
    defer signal.Stop(sigChannel)

    go func() {
		sig := <-sigChannel
		logger.Infof("Received signal: %s, initiating shutdown...", sig)
		cancel()
	}()

    // repo := postgresql.NewPostgreSQLConnection(cfg.Database, shutdownChan)
	// logger.Info("Postgresql conn has been established")
    
    // cache := valkey.NewValkeyClient(cfg.Valkey, shutdownChan)
	// logger.Info("Valkey conn has been established")

    wg := new(sync.WaitGroup)
    wg.Add(1)

    app := app.NewBU(logger, wg, ctx)

    go app.RunApplication()
	<-ctx.Done()
	logger.Info("Context cancelled, shutting down...")

	app.Shutdown()

    wg.Wait()
    logger.Info("App exited cleanly.")
}