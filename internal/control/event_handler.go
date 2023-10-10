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
	usr, err = eh.toInnerUser(upd.SentFrom())
	if err != nil {
		return
	}

	// Getting inner TgChat
	var tgChat *tcPkg.TgChat
	tgChat, err = eh.toInnerTgChat(usr, upd.FromChat())
	if err != nil {
		return
	}

	// Ignore non-message event
	if upd.Message == nil {
		return
	}

	msg := upd.Message

	err = eh.processMessage(tgChat, msg)
	if err != nil {
		return
	}

	err = eh.sendReplyMessage(tgChat, msg)
	if err != nil {
		return
	}
}

func (eh EventHandler) toInnerUser(outerU *tgbotapi.User) (u *usrPkg.User, err error) {
	u = usrPkg.MapToInner(outerU)

	var userExists bool
	if userExists, err = eh.Repos.User.IsExists(u); err == nil {
		if userExists {
			err = eh.Repos.User.Update(u)
		} else {
			err = eh.Repos.User.SaveNew(u)
		}
	}

	return u, err
}

func (eh EventHandler) toInnerTgChat(u *usrPkg.User, outerChat *tgbotapi.Chat) (tc *tcPkg.TgChat, err error) {
	if tc, err = eh.Repos.TgChat.TgChatByUserId(u.Id); err == nil && tc == nil {
		state, _ := eh.Repos.TgChat.StartState()
		tc = tcPkg.NewTgChat(outerChat, u.Id, state)
		err = eh.Repos.TgChat.SaveNewTgChat(tc)
	}

	return tc, err
}

func (eh EventHandler) processMessage(tc *tcPkg.TgChat, msg *tgbotapi.Message) error {
	// Message validation
	var err error
	err = tc.ValidateMessage(msg)
	if err != nil {
		return err
	}

	//	Отразить требуемые командой измения в польз. данных

	// Setting tgChat next state
	var cmd *tcPkg.Cmd
	cmd, err = eh.Repos.TgChat.CmdByCode(msg.CommandWithAt())
	if err != nil {
		return err
	}
	tc.SetState(cmd.DestState)
	err = eh.Repos.TgChat.UpdateTgChatState(tc)
	if err != nil {
		return err
	}

	return nil
}

func (eh EventHandler) sendReplyMessage(tc *tcPkg.TgChat, iMsg *tgbotapi.Message) error {
	oMsgText := tc.State.Msg
	if tc.State.MsgHdr != "" {
		oMsgText = tc.State.MsgHdr + "\n\n" + oMsgText
	}
	if tc.State.MsgFtr != "" {
		oMsgText += "\n\n" + tc.State.MsgFtr
	}

	msg := tgbotapi.NewMessage(iMsg.Chat.ID, oMsgText)
	// msg.ReplyToMessageID = iMsg.MessageID

	_, err := eh.TgBotAPI.Send(msg)

	return err
}
