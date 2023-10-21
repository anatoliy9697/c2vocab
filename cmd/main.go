package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/anatoliy9697/c2vocab/internal/control"
	tcRepo "github.com/anatoliy9697/c2vocab/internal/model/tgchat/repo"
	usrRepo "github.com/anatoliy9697/c2vocab/internal/model/user/repo"
	wlRepo "github.com/anatoliy9697/c2vocab/internal/model/wordlist/repo"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var err error

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})) // TODO: Вынести определение уровня логирования в параметры окружения или в конфиг

	defer func() {
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	logger.Info("C2Vocab initialization")

	// Creating subsidiary context and assigning it to external interruptions listening
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Storage connection pool initialization
	var pgPool *pgxpool.Pool
	if pgPool, err = pgxpool.New(mainCtx, os.Getenv("POSTGRES_CONN_STRING")); err != nil {
		return
	}
	defer pgPool.Close()

	// Telegram client initialization
	var tgBotAPI *tgbotapi.BotAPI
	if tgBotAPI, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN")); err != nil {
		return
	}
	// tgBotAPI.Debug = true

	// Event fetcher (ef) configurating and running
	var tgchatRepo tcRepo.Repo
	if tgchatRepo, err = tcRepo.Init(mainCtx, pgPool); err != nil {
		return
	}
	res := res.Resources{
		UsrRepo:  usrRepo.Init(mainCtx, pgPool),
		TcRepo:   tgchatRepo,
		WLRepo:   wlRepo.Init(mainCtx, pgPool),
		TgBotAPI: tgBotAPI,
		Logger:   logger,
	}
	efDone := make(chan struct{})
	go control.EventFetcher{
		TgBotUpdsOffset:       0,
		TgBotUpdsTimeout:      30,
		MaxEventHandlers:      10,
		WaitForHandlerTimeout: 100,
		Res:                   res,
	}.Run(mainCtx, efDone)

	logger.Info("C2Vocab is running")

	// Keeping alive
	<-mainCtx.Done()

	logger.Info("Shutdown initialized. Waiting for all subsidiary goroutines finishing")

	// Waiting for event fetcher completion
	<-efDone

	logger.Info("C2Vocab execution completed")
}
