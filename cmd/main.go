package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/anatoliy9697/c2vocab/internal/control"
	tskRepo "github.com/anatoliy9697/c2vocab/internal/model/task/repo"
	tcRepo "github.com/anatoliy9697/c2vocab/internal/model/tgchat/repo"
	usrRepo "github.com/anatoliy9697/c2vocab/internal/model/user/repo"
	wlRepo "github.com/anatoliy9697/c2vocab/internal/model/wordlist/repo"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	var err error

	viper.SetConfigFile("./config/config.yaml")
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("reading config file fatal error: " + err.Error())
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(viper.GetInt("log_level"))}))

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

	// Other resources initialization
	var tgchatRepo tcRepo.Repo
	if tgchatRepo, err = tcRepo.Init(mainCtx, pgPool); err != nil {
		return
	}
	res := res.Resources{
		UsrRepo:  usrRepo.Init(mainCtx, pgPool),
		TcRepo:   tgchatRepo,
		WLRepo:   wlRepo.Init(mainCtx, pgPool),
		TskRepo:  tskRepo.Init(mainCtx, pgPool),
		TgBotAPI: tgBotAPI,
		Logger:   logger,
	}
	goroutinesDone := make(chan struct{}, 2)

	// Event fetcher (ef) configurating and running
	go control.EventFetcher{
		TgBotUpdsOffset:  viper.GetInt("tg_bot_upds_offset"),
		TgBotUpdsTimeout: viper.GetInt("tg_bot_upds_timeout"),
		MaxEventHandlers: viper.GetInt("max_event_handlers"),
		Res:              res,
	}.Run(mainCtx, goroutinesDone)

	// Scheduler configurating and running
	go control.Scheduler{
		MaxTaskHandlers:    viper.GetInt("max_tack_handlers"),
		TaskWaitingTime:    viper.GetInt("task_waiting_time"),
		TaskBatchSize:      viper.GetInt("task_batch_size"),
		MaxTimeForReassign: viper.GetInt("max_time_for_reassign"),
		Res:                res,
	}.Run(mainCtx, goroutinesDone)

	logger.Info("C2Vocab is running")

	// Keeping alive
	<-mainCtx.Done()

	logger.Info("Shutdown initialized. Waiting for all subsidiary goroutines finishing")

	// Waiting for event fetcher and scheduler completion
	<-goroutinesDone
	<-goroutinesDone

	logger.Info("C2Vocab execution completed")
}
