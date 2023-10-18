package tgchat

import (
	"errors"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgChat struct {
	TgId         int64
	UserId       int32
	CreatedAt    time.Time
	State        *State
	WLFrgnLang   *commons.Lang
	WLNtvLang    *commons.Lang
	WLId         int32
	WL           *wlPkg.WordList
	BotLastMsgId int
}

type State struct {
	Code              string
	MsgHdr            string
	MsgBody           string
	MsgFtr            string
	WaitForWLFrgnLang bool
	WaitForWLNtvLang  bool
	WaitForWLName     bool
	StateCmd          *Cmd
	NextStateCode     string
	AvailCmds         [][]*Cmd
}

type Cmd struct {
	Code          string
	DisplayLabel  string
	DestStateCode string
}

type IncomingMsg struct {
	Id      int
	Cmd     *Cmd
	CmdArgs []string
	Text    string
}

var (
	ErrEmptyCmd            = errors.New("получена пустая команда")
	ErrUnexpectedCmd       = errors.New("получена неожиданная команда")
	ErrUnexpectedDataInput = errors.New("ожидается команда, не ввод данных")
	ErrEmptyDataInput      = errors.New("получена пустая строка в качестве входных данных")
)

func (tc *TgChat) SetState(s *State) {
	tc.State = s
}

func (tc *TgChat) SetBotLastMsgId(msgId int) {
	tc.BotLastMsgId = msgId
}

func (tc *TgChat) TgOutgoingMsgText() string {
	msgText := tc.State.MsgBody
	if tc.State.MsgHdr != "" {
		msgText = tc.State.MsgHdr + "\n\n" + msgText
	}
	if tc.State.MsgFtr != "" {
		msgText += "\n\n" + tc.State.MsgFtr
	}

	return msgText
}

func (s State) IsWaitForDataInput() bool {
	return s.WaitForWLName
}

func (s State) IsCmdAvail(cmdCode string) bool {
	for _, cmdsRow := range s.AvailCmds {
		for _, cmd := range cmdsRow {
			if cmd.Code == cmdCode {
				return true
			}
		}
	}

	if s.StateCmd != nil && s.StateCmd.Code == cmdCode {
		return true
	}

	return false
}

func (s State) TgInlineKeyboradAvailCmdsRows() [][]tgbotapi.InlineKeyboardButton {
	inlineKeyboardRows := make([][]tgbotapi.InlineKeyboardButton, len(s.AvailCmds))
	for i, cmdsRow := range s.AvailCmds {
		inlineKeyboardRows[i] = make([]tgbotapi.InlineKeyboardButton, len(cmdsRow))
		for j, cmd := range cmdsRow {
			inlineKeyboardRows[i][j] = cmd.TgButton()
		}
	}

	return inlineKeyboardRows
}

func (c Cmd) TgButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(c.DisplayLabel, c.Code)
}
