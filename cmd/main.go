package main

import (
	"context"
	app "horsey/internal/application"
	"horsey/internal/config"
	"horsey/internal/handlers/tgBot"
	"horsey/internal/router/tgBot"
	"horsey/internal/storage"
	"horsey/internal/usecase"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"gopkg.in/telebot.v4"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	db, err := storage.New(ctx, cfg.Postgres.Addr)
	if err != nil {
		log.Error("failed to connect to storage", "error", err)
		return
	}
	defer db.Close()

	dbPersistence, err := storage.New(ctx, cfg.Postgres.Addr)
	if err != nil {
		log.Error("failed to connect to storage", "error", err)
		return
	}

	useCase := usecase.NewUseCase(log, dbPersistence)

	bot, err := botRun(cfg, log, useCase)
	if err != nil {
		log.Error("Failed to start botRun", "err", err)
		os.Exit(1)
	}

	router.Router(ctx, bot)

	application := app.NewApp(ctx, cfg, log, bot)

	application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	application.Shutdown()
}

func botRun(cfg *config.Config, log *slog.Logger, useCase *usecase.UseCase) (*tgBot.TgBot, error) {
	pref := telebot.Settings{
		Token: cfg.TgBot.TgToken,
		Poller: &telebot.LongPoller{
			Timeout:      10 * time.Second,
			LastUpdateID: -1,
		},
	}

	newBot, err := telebot.NewBot(pref)
	if err != nil {
		log.Error("BotRun: failed creating a bot with error", err.Error())
		return nil, err
	}

	bot := tgBot.New(log, newBot, useCase)

	return bot, nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
