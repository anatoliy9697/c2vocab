package control

import (
	"fmt"
	"strings"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
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
	var tc *tcPkg.TgChat
	tc, err = eh.toInnerTgChat(usr, upd.FromChat())
	if err != nil {
		return
	}

	// Ignore non-message and non-command events
	if upd.Message == nil && upd.CallbackQuery == nil {
		return
	}

	err = eh.processUpdate(tc, upd)
	if err != nil {
		return
	}

	err = eh.sendReplyMessage(tc)
	if err != nil {
		return
	}

	err = eh.Repos.TgChat.UpdateTgChat(tc)
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
		tc = &tcPkg.TgChat{
			TgId:      outerChat.ID,
			UserId:    u.Id,
			State:     state,
			CreatedAt: time.Now(),
		}
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
		if !tc.State.IsCmdAvail(cmdCode) {
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

func (eh EventHandler) processIncomingMsg(tc *tcPkg.TgChat, msg *tcPkg.IncomingMsg) error {
	var err error

	switch {
	case msg.Cmd != nil && msg.Cmd.Code == "wl_frgn_lang":
		tc.WLFrgnLang = commons.LangByCode(msg.CmdArgs[0])
	case msg.Cmd != nil && msg.Cmd.Code == "wl_ntv_lang":
		tc.WLNtvLang = commons.LangByCode(msg.CmdArgs[0])
	case msg.Text != "" && tc.State.WaitForWLName:
		wl := &wlPkg.WordList{
			Active:    true,
			Name:      msg.Text,
			FrgnLang:  tc.WLFrgnLang,
			NtvLang:   tc.WLNtvLang,
			OwnerId:   tc.UserId,
			CreatedAt: time.Now(),
		}
		err = eh.Repos.WL.SaveNewWL(wl)
	}

	return err
}

func (eh EventHandler) processUpdate(tc *tcPkg.TgChat, upd *tgbotapi.Update) error {
	var err error
	var msg *tcPkg.IncomingMsg

	msg, err = eh.validateAndMapToIncMsg(tc, upd)
	if err != nil {
		return err
	}

	err = eh.processIncomingMsg(tc, msg)
	if err != nil {
		return err
	}

	// Set tgChat next state
	var nextState *tcPkg.State
	if msg.Cmd != nil {
		nextState, err = eh.Repos.TgChat.StateByCode(msg.Cmd.DestStateCode)
	} else if tc.State.NextStateCode != "" {
		nextState, err = eh.Repos.TgChat.StateByCode(tc.State.NextStateCode)
	}
	if err != nil {
		return err
	}
	tc.SetState(nextState)

	return nil
}

func (eh EventHandler) sendReplyMessage(tc *tcPkg.TgChat) error {
	_, err := eh.TgBotAPI.Send(tc.TgOutgoingMsg())

	return err
}
