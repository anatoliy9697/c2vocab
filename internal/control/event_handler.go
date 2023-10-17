package control

import (
	"fmt"
	"strconv"
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
	if usr, err = eh.toInnerUser(upd.SentFrom()); err != nil {
		return
	}

	// Getting inner TgChat
	var tc *tcPkg.TgChat
	if tc, err = eh.toInnerTgChat(usr, upd.FromChat()); err != nil {
		return
	}

	// Ignore non-message and non-command events
	if upd.Message == nil && upd.CallbackQuery == nil {
		return
	}

	if err = eh.processUpdate(tc, upd); err != nil {
		return
	}

	if err = eh.sendReplyMessage(tc); err != nil {
		return
	}

	if err = eh.Repos.TgChat.UpdateTgChat(tc); err != nil {
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
		if cmd, err = eh.Repos.TgChat.CmdByCode(cmdCode); err != nil {
			return nil, err
		}
	}

	return &tcPkg.IncomingMsg{
		Cmd:     cmd,
		CmdArgs: cmdArgs,
		Text:    msgText,
	}, nil
}

func (eh EventHandler) processIncomingMsg(tc *tcPkg.TgChat, msg *tcPkg.IncomingMsg) error {
	var err error

	switch {
	case msg.Cmd != nil && (msg.Cmd.Code == "start" || msg.Cmd.Code == "to_main_menu"):
		tc.WLFrgnLang = nil
		tc.WLNtvLang = nil
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
		if err = eh.Repos.WL.SaveNewWL(wl); err != nil {
			return err
		}
		tc.WLFrgnLang = nil
		tc.WLNtvLang = nil
	}

	return nil
}

func (eh EventHandler) processUpdate(tc *tcPkg.TgChat, upd *tgbotapi.Update) error {
	var err error
	var msg *tcPkg.IncomingMsg

	if msg, err = eh.validateAndMapToIncMsg(tc, upd); err != nil {
		return err
	}

	if err = eh.processIncomingMsg(tc, msg); err != nil {
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

func (eh EventHandler) tgInlineKeyboradStateCmdRows(tc *tcPkg.TgChat) ([][]tgbotapi.InlineKeyboardButton, error) {
	var err error
	var inlineKeyboardRows [][]tgbotapi.InlineKeyboardButton

	switch { // TODO: Повыносить все case'ы нахрен в функции
	case tc.State.WaitForWLFrgnLang:
		rowsCount := len(commons.AvailLangs)
		inlineKeyboardRows = make([][]tgbotapi.InlineKeyboardButton, rowsCount)
		for i, lang := range commons.AvailLangs {
			inlineKeyboardRows[i] = make([]tgbotapi.InlineKeyboardButton, 1)
			inlineKeyboardRows[i][0] = tgbotapi.NewInlineKeyboardButtonData(lang.Name, tc.State.StateCmd.Code+" "+lang.Code)
		}
	case tc.State.WaitForWLNtvLang:
		rowsCount := len(commons.AvailLangs) - 1
		inlineKeyboardRows = make([][]tgbotapi.InlineKeyboardButton, rowsCount)
		i := 0
		for _, lang := range commons.AvailLangs {
			if lang.Code != tc.WLFrgnLang.Code {
				inlineKeyboardRows[i] = make([]tgbotapi.InlineKeyboardButton, 1)
				inlineKeyboardRows[i][0] = tgbotapi.NewInlineKeyboardButtonData(lang.Name, tc.State.StateCmd.Code+" "+lang.Code)
				i++
			}
		}
	case tc.State.StateCmd != nil && tc.State.StateCmd.Code == "wl":
		var wls []*wlPkg.WordList
		if wls, err = eh.Repos.WL.ActiveWLByOwnerId(tc.UserId); err != nil {
			return nil, err
		}
		var inlineKeyboardRow []tgbotapi.InlineKeyboardButton
		for _, wl := range wls {
			inlineKeyboardRow = []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(wl.Name, tc.State.StateCmd.Code+" "+strconv.Itoa(int(wl.Id)))}
			inlineKeyboardRows = append(inlineKeyboardRows, inlineKeyboardRow)
		}
	}

	return inlineKeyboardRows, nil
}

func (eh EventHandler) tgInlineKeyboradMarkup(tc *tcPkg.TgChat) (tgbotapi.InlineKeyboardMarkup, error) {
	var inlineKeyboardMarkup tgbotapi.InlineKeyboardMarkup
	inlineKeyboardRows := make([][]tgbotapi.InlineKeyboardButton, 0)

	inlineKeyboardStateCmdRows, err := eh.tgInlineKeyboradStateCmdRows(tc)
	if err != nil {
		return inlineKeyboardMarkup, err
	}
	if len(inlineKeyboardStateCmdRows) > 0 {
		inlineKeyboardRows = append(inlineKeyboardRows, inlineKeyboardStateCmdRows...)
	}

	inlineKeyboardAvailCmdsRows := tc.State.TgInlineKeyboradAvailCmdsRows()
	if len(inlineKeyboardAvailCmdsRows) > 0 {
		inlineKeyboardRows = append(inlineKeyboardRows, inlineKeyboardAvailCmdsRows...)
	}

	if len(inlineKeyboardRows) > 0 {
		inlineKeyboardMarkup = tgbotapi.NewInlineKeyboardMarkup(inlineKeyboardRows...)
	}

	return inlineKeyboardMarkup, nil
}

func (eh EventHandler) tgOutgoingMsg(tc *tcPkg.TgChat) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(tc.TgId, tc.TgOutgoingMsgText())

	// TgChat control buttons
	replyMarkup, err := eh.tgInlineKeyboradMarkup(tc)
	if err != nil {
		return msg, err
	}

	msg.ReplyMarkup = replyMarkup

	return msg, nil
}

func (eh EventHandler) tgMsgEditing(tc *tcPkg.TgChat) (tgbotapi.EditMessageTextConfig, error) {
	var editMsgConfig tgbotapi.EditMessageTextConfig

	// TgChat control buttons
	replyMarkup, err := eh.tgInlineKeyboradMarkup(tc)
	if err != nil {
		return editMsgConfig, err
	}

	editMsgConfig = tgbotapi.NewEditMessageTextAndMarkup(
		tc.TgId,
		tc.BotLastMsgId,
		tc.TgOutgoingMsgText(),
		replyMarkup,
	)

	return editMsgConfig, nil
}

func (eh EventHandler) sendReplyMessage(tc *tcPkg.TgChat) error {
	var err error
	var msg tgbotapi.Chattable

	if tc.BotLastMsgId != 0 {
		msg, err = eh.tgMsgEditing(tc)
	} else {
		msg, err = eh.tgOutgoingMsg(tc)
	}
	if err != nil {
		return err
	}

	var msgInTg tgbotapi.Message
	if msgInTg, err = eh.TgBotAPI.Send(msg); err != nil && tc.BotLastMsgId != 0 {
		if msg, err = eh.tgOutgoingMsg(tc); err != nil {
			return err
		}
		if msgInTg, err = eh.TgBotAPI.Send(msg); err != nil {
			return err
		}
	}

	tc.SetBotLastMsgId(msgInTg.MessageID)

	return nil
}
