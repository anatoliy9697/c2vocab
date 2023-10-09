package control

import (
	"fmt"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EventHandler struct {
	HandlerCode string
	TgBotAPI    *tgbotapi.BotAPI
	Repos       Repos
}

func (eh EventHandler) Run(done chan string, upd tgbotapi.Update) {
	defer func() { done <- eh.HandlerCode }()
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err.Error()) // TODO: Сделать адекватное логирование + сообщение об ошибке пользователю
		}
	}()

	// Getting inner user
	var usr *usrPkg.User
	usr, err = eh.Repos.User.ToInner(upd.SentFrom())
	if err != nil {
		return
	}

	// Getting inner TgChat
	var tgChat *tcPkg.TgChat
	tgChat, err = eh.Repos.TgChat.ToInnerTgChat(usr, upd.FromChat())
	if err != nil {
		return
	}

	// Ignore non-message event
	if upd.Message == nil {
		return
	}

	msg := upd.Message

	err = tgChat.ValidateMessage(msg)
	if err != nil {
		return
	}

	//	Отразить требуемые командой измения в польз. данных

	//	Сменить состояние чата

	//	Сформировать сообщение пользователю согласно текущему состоянию и отправить его в чат

	if upd.Message != nil {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
		msg.ReplyToMessageID = upd.Message.MessageID
		eh.TgBotAPI.Send(msg)
	}
}
