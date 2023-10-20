package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anatoliy9697/c2vocab/internal/control"
	tcRepo "github.com/anatoliy9697/c2vocab/internal/model/tgchat/repo"
	usrRepo "github.com/anatoliy9697/c2vocab/internal/model/user/repo"
	wlRepo "github.com/anatoliy9697/c2vocab/internal/model/wordlist/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Creating subsidiary context and assigning it to external interruptions listening
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Storage connection pool initialization
	pgPool, err := pgxpool.New(mainCtx, os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		log.Fatal(err) // TODO: Сделать адекватное логирование и завершение программы
	}
	defer pgPool.Close()

	// Telegram client initialization
	tgBotAPI, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err) // TODO: Сделать адекватное логирование и завершение программы
	}
	// tgBotAPI.Debug = true

	// Event fetcher (ef) configurating and running
	tgchatRepo, err := tcRepo.Init(mainCtx, pgPool)
	if err != nil {
		log.Fatal(err) // TODO: Сделать адекватное логирование и завершение программы
	}
	efDone := make(chan struct{})
	go control.EventFetcher{
		TgBotAPI:              tgBotAPI,
		TgBotUpdsOffset:       0,
		TgBotUpdsTimeout:      30,
		MaxEventHandlers:      10,
		WaitForHandlerTimeout: 100,
		Repos: control.Repos{
			User:   usrRepo.Init(mainCtx, pgPool),
			TgChat: tgchatRepo,
			WL:     wlRepo.Init(mainCtx, pgPool),
		},
	}.Run(mainCtx, efDone)

	// Keeping alive
	<-mainCtx.Done()

	fmt.Println("Производится grasefull shutdown. Ждем завершения дочерних горутин") // TODO: Сделать адекватное логирование

	// Waiting for event fetcher completion
	<-efDone

	fmt.Println("Все горутингы завершили свою работу") // TODO: Сделать адекватное логирование
}
