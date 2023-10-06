package control

import (
	"fmt"

	tgchatPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EventHandler struct {
	HandlerCode string
	TgBotAPI    *tgbotapi.BotAPI
	Repos       *Repos
}

func (eh EventHandler) Run(done chan string, upd *tgbotapi.Update) {
	defer func() { done <- eh.HandlerCode }()
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err.Error()) // TODO: Сделать адекватное логирование + сообщение об ошибке пользователю
		}
	}()

	outerUsr := upd.SentFrom()
	usr := usrPkg.MapToInner(outerUsr)

	// User checking and saving/updating
	var newUser bool
	newUser, err = eh.Repos.User.Set(usr)
	if err != nil {
		return
	}

	// Getting chat and state
	var state *tgchatPkg.State
	var chat *tgchatPkg.TgChat
	if newUser {
		state, _ = eh.Repos.TgChat.StartState()
		chat = tgchatPkg.MapToInner(upd.FromChat())
		chat.StateCode = state.Code
		chat.UserId = usr.Id
		err = eh.Repos.TgChat.SaveNew(chat)
		if err != nil {
			return
		}
	} else {
		//		получить текущее состояние чата с пользователем
		chat, err = eh.Repos.TgChat.TgChatByUserId(usr.Id)
		if err != nil {
			return
		}
		state, err = eh.Repos.TgChat.StateByCode(chat.StateCode)
		if err != nil {
			return
		}
	}

	//	Получить список доступных для текущего состояния действий

	//	Если выполненное пользователем действие не соответствует доступным, то
	//		сформировать сообщение о недопустимом действии, отправить пользователю и завершить свое выполнение

	//	Отразить требуемые командой измения в польз. данных

	//	Сменить состояние чата

	//	Сформировать сообщение пользователю согласно текущему состоянию и отправить его в чат

	if upd.Message != nil {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
		msg.ReplyToMessageID = upd.Message.MessageID
		eh.TgBotAPI.Send(msg)
	}
}
