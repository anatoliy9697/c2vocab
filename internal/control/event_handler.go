package control

import (
	"fmt"
	"strings"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	"github.com/anatoliy9697/c2vocab/pkg/sliceutils"
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

	// Ignore non-message and non-command events
	if upd.Message == nil && upd.CallbackQuery == nil {
		return
	}

	err = eh.processUpdate(tgChat, upd)
	if err != nil {
		return
	}

	err = eh.sendReplyMessage(tgChat)
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
		tc = tcPkg.NewTgChat(outerChat.ID, u.Id, state)
		err = eh.Repos.TgChat.SaveNewTgChat(tc)
	}

	return tc, err
}

func (eh EventHandler) validateAndMapToIncMsg(tc *tcPkg.TgChat, upd *tgbotapi.Update) (*tcPkg.IncomingMsg, error) {
	var err error
	var cmdCode, msgText string
	var cmdArgs []string
	var cmd *tcPkg.Cmd

	if upd.CallbackQuery != nil {
		if upd.CallbackQuery.Data == "" {
			return nil, tcPkg.ErrEmptyCmd
		}
		cmdParts := strings.Split(upd.CallbackQuery.Data, " ")
		if len(cmdParts) > 1 {
			cmdArgs = cmdParts[1:]
		}
		cmdCode = strings.ReplaceAll(cmdParts[0], "/", "")
		if !sliceutils.IsStrInSlice(tc.State.AvailCmdCodes, cmdCode) {
			return nil, tcPkg.ErrUnexpectedCmd
		}
	} else if upd.Message != nil {
		if upd.Message.IsCommand() {
			cmdCode = upd.Message.Command()
			cmdArgs = strings.Split(upd.Message.CommandArguments(), " ")
			if cmdCode != "start" {
				return nil, tcPkg.ErrUnexpectedCmd
			}
		} else {
			msgText = upd.Message.Text
			if !tc.State.IsWaitForDataInput() {
				return nil, tcPkg.ErrUnexpectedDataInput
			}
		}
	}

	if cmdCode != "" {
		cmd, err = eh.Repos.TgChat.CmdByCode(cmdCode)
	}
	if err != nil {
		return nil, err
	}

	return tcPkg.NewIncomingMsg(cmd, cmdArgs, msgText), nil
}

func (eh EventHandler) processUpdate(tc *tcPkg.TgChat, upd *tgbotapi.Update) error {
	var err error
	var msg *tcPkg.IncomingMsg

	msg, err = eh.validateAndMapToIncMsg(tc, upd)
	if err != nil {
		return err
	}

	// Тут будет обработка данных пользователя

	tc.SetState(msg.Cmd.DestState)

	err = eh.Repos.TgChat.UpdateTgChat(tc)
	if err != nil {
		return err
	}

	return nil
}

func (eh EventHandler) sendReplyMessage(tc *tcPkg.TgChat) error {
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("На старт", "/to_start"),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "/start"),
		),
	)

	oMsgText := tc.State.Msg
	if tc.State.MsgHdr != "" {
		oMsgText = tc.State.MsgHdr + "\n\n" + oMsgText
	}
	if tc.State.MsgFtr != "" {
		oMsgText += "\n\n" + tc.State.MsgFtr
	}

	msg := tgbotapi.NewMessage(tc.TgId, oMsgText)
	// msg.ReplyToMessageID = iMsg.MessageID

	msg.ReplyMarkup = keyboard

	_, err := eh.TgBotAPI.Send(msg)

	return err
}
