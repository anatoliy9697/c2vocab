package control

import (
	"fmt"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	"github.com/anatoliy9697/c2vocab/internal/usecases"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EventHandler struct {
	HandlerCode string
	TgBotAPI    *tgbotapi.BotAPI
	Repos       Repos
}

func (eh EventHandler) Run(done chan string, upd *tgbotapi.Update) {
	defer func() { done <- eh.HandlerCode }()

	var err error

	defer func() {
		if err != nil {
			fmt.Println(err.Error()) // TODO: Сделать адекватное логирование + сообщение об ошибке пользователю
		}
	}()

	var usr *usrPkg.User
	if usr, err = usecases.MapToInnerUserAndSave(eh.Repos.User, upd.SentFrom()); err != nil {
		return
	}

	var tc *tcPkg.Chat
	if tc, err = usecases.MapToInnerTgChatAndSave(eh.Repos.TgChat, eh.Repos.WL, upd.FromChat(), usr); err != nil {
		return
	}

	// Ignore non-message and non-command events
	if upd.Message == nil && upd.CallbackQuery == nil {
		return
	}

	if err = usecases.ProcessUpd(eh.Repos.TgChat, eh.Repos.WL, eh.TgBotAPI, tc, upd); err != nil {
		return
	}

	if err = usecases.SendReplyMessage(eh.Repos.WL, eh.TgBotAPI, tc); err != nil {
		return
	}

	if err = eh.Repos.TgChat.UpdateTgChat(tc); err != nil {
		return
	}
}
