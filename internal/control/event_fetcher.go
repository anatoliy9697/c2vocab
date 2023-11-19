package control

import (
	"context"

	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type EventFetcher struct {
	TgBotUpdsOffset  int
	TgBotUpdsTimeout int
	MaxEventHandlers int
	Res              res.Resources
}

func (ef EventFetcher) Run(ctx context.Context, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	ef.Res.Logger.Info("Event fetcher is running")

	updConfig := tgbotapi.NewUpdate(ef.TgBotUpdsOffset)
	updConfig.Timeout = ef.TgBotUpdsTimeout
	upds := ef.Res.TgBotAPI.GetUpdatesChan(updConfig)

	handlers := make(map[string]struct{}, ef.MaxEventHandlers)
	handlerDone := make(chan string, 10)
	handlerCode := ""

loop:
	for {
		select {

		// Event fetcher shutdown
		case <-ctx.Done():
			break loop

		// Handler had finished
		case handlerCode = <-handlerDone:
			delete(handlers, handlerCode)

		// Got new update
		case upd := <-upds:
			if len(handlers) >= ef.MaxEventHandlers {
				ef.Res.Logger.Info("No free event handlers. Waiting for handler")
				handlerCode = <-handlerDone
				delete(handlers, handlerCode)
			}
			handlerCode := uuid.NewString()[:7]
			handlers[handlerCode] = struct{}{}
			ef.Res.Logger.Info("Running event handler " + handlerCode)
			go EventHandler{
				Code: handlerCode,
				Res:  ef.Res,
			}.Run(handlerDone, &upd)

		}
	}

	// Whaiting for rest handlers finishing
	for len(handlers) > 0 {
		handlerCode := <-handlerDone
		delete(handlers, handlerCode)
	}

	ef.Res.Logger.Info("Event fetcher execution completed")
}
