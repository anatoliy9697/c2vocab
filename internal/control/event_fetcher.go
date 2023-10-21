package control

import (
	"context"
	"time"

	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type EventFetcher struct {
	TgBotUpdsOffset       int
	TgBotUpdsTimeout      int
	MaxEventHandlers      int
	WaitForHandlerTimeout int // ms
	Res                   res.Resources
}

func (ef EventFetcher) Run(ctx context.Context, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	ef.Res.Logger.Info("Event-fetcher is running")

	updConfig := tgbotapi.NewUpdate(ef.TgBotUpdsOffset)
	updConfig.Timeout = ef.TgBotUpdsTimeout
	upds := ef.Res.TgBotAPI.GetUpdatesChan(updConfig)

	handlers := make(map[string]struct{}, ef.MaxEventHandlers)
	handlerDone := make(chan string, 10)

outer:
	for {
		select {

		// Event fetcher shutdown
		case <-ctx.Done():
			break outer

		// Handler had finished
		case handlerCode := <-handlerDone:
			delete(handlers, handlerCode)

		// Got new update
		case upd := <-upds:
		inner:
			for {
				if len(handlers) < ef.MaxEventHandlers {
					handlerCode := uuid.NewString()[:7]
					handlers[handlerCode] = struct{}{}
					ef.Res.Logger.Info("Running event-handler " + handlerCode)
					go EventHandler{
						Code: handlerCode,
						Res:  ef.Res,
					}.Run(handlerDone, &upd)
					break inner
				} else {
					ef.Res.Logger.Info("No free event-handlers. Waiting for timeout")
					time.Sleep(time.Duration(ef.WaitForHandlerTimeout) * time.Millisecond)
				}
			}

		}
	}

	// Whaiting for rest handlers finishing
	for len(handlers) > 0 {
		handlerCode := <-handlerDone
		delete(handlers, handlerCode)
	}

	ef.Res.Logger.Info("Event-fetcher execution completed")
}
