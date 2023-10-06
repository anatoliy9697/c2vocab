package control

import (
	"context"
	"fmt"
	"time"

	tgchat "github.com/anatoliy9697/c2vocab/internal/model/tgchat/repo"
	usr "github.com/anatoliy9697/c2vocab/internal/model/user/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type Repos struct {
	User   usr.Repo
	TgChat tgchat.Repo
}

type EventFetcher struct {
	TgBotAPI              *tgbotapi.BotAPI
	TgBotUpdsOffset       int
	TgBotUpdsTimeout      int
	MaxEventHandlers      int
	WaitForHandlerTimeout int // ms
	Repos                 *Repos
}

func (ef EventFetcher) Run(ctx context.Context, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	fmt.Println("Event fetcher начал свою работу")

	updConfig := tgbotapi.NewUpdate(ef.TgBotUpdsOffset)
	updConfig.Timeout = ef.TgBotUpdsTimeout
	upds := ef.TgBotAPI.GetUpdatesChan(updConfig)

	handlers := make(map[string]struct{}, ef.MaxEventHandlers)
	handlerDone := make(chan string, 10)

loop:
	for {
		select {

		// Event fetcher shutdown
		case <-ctx.Done():
			break loop

		// Handler had finished
		case handlerCode := <-handlerDone:
			delete(handlers, handlerCode)

		// Got new update
		case upd := <-upds:
			if len(handlers) < ef.MaxEventHandlers {
				handlerCode := uuid.NewString()
				handlers[handlerCode] = struct{}{}
				go EventHandler{
					HandlerCode: handlerCode,
					TgBotAPI:    ef.TgBotAPI,
					Repos:       ef.Repos,
				}.Run(handlerDone, &upd)
			} else {
				time.Sleep(time.Duration(ef.WaitForHandlerTimeout) * time.Millisecond)
			}

		}
	}

	// Whaiting for rest handlers finishing
	for len(handlers) > 0 {
		handlerCode := <-handlerDone
		delete(handlers, handlerCode)
	}

	fmt.Println("Event fetcher завершил свою работу") // TODO: Сделать адекватное логирование
}
