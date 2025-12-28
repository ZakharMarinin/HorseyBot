package app

import (
	"context"
	"horsey/internal/config"
	"horsey/internal/handlers/tgBot"
	"log/slog"
	"sync"
)

type App struct {
	ctx   context.Context
	cfg   *config.Config
	log   *slog.Logger
	tgbot *tgBot.TgBot
}

func NewApp(ctx context.Context, cfg *config.Config, log *slog.Logger, bot *tgBot.TgBot) *App {
	return &App{
		ctx:   ctx,
		cfg:   cfg,
		log:   log,
		tgbot: bot,
	}
}

func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.tgbot.BackgroundTimer(a.ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.log.Info("Run: bot starting...")
		a.tgbot.Bot.Start()
	}()

	a.log.Info("Run: bot and worker are successfully initiated")

	go func() {
		wg.Wait()
		a.log.Info("Run: bot stopped")
	}()

	return nil
}

func (a *App) Shutdown() {
	a.log.Info("Shutdown")

	a.tgbot.Bot.Stop()
}
