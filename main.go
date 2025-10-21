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

	"github.com/ayayaakasvin/trends-updater/internal/repo/postgresql"
	"github.com/ayayaakasvin/trends-updater/internal/repo/valkey"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    cfg := config.MustLoadConfig()
    logger := logger.SetupLogger()

	sigChannel := make(chan os.Signal, 1)
    signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)
    defer signal.Stop(sigChannel)

    SetupChan := inner.NewShutdownChannel()
    go func() {
        select {
        case sig := <-sigChannel:
            logger.Infof("Received signal: %s", sig)
        case msg := <-SetupChan:
            logger.Errorf("Internal error: %s", msg)
        }
        cancel()
    }()

    er := postgresql.NewPostgreSQLConnection(cfg.Database, SetupChan)
    cc := valkey.NewValkeyClient(cfg.Valkey, SetupChan)

    wg := new(sync.WaitGroup)
    wg.Add(1)

    app := app.NewBU(logger, wg, ctx, er, cc)

    go app.RunApplication()

	<-ctx.Done()
	logger.Info("Context cancelled, shutting down...")

    wg.Wait()
    logger.Info("App exited cleanly.")
}