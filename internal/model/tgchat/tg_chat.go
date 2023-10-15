package tgchat

import (
	"errors"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgChat struct {
	TgId       int64
	UserId     int32
	State      *State
	WLFrgnLang *commons.Lang
	WLNtvLang  *commons.Lang
	CreatedAt  time.Time
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

func NewTgChat(tgId int64, userId int32, state *State) *TgChat {
	return &TgChat{
		TgId:   tgId,
		UserId: userId,
		State:  state,
	}
}

func NewIncomingMsg(cmd *Cmd, cmdArgs []string, msgText string) *IncomingMsg {
	return &IncomingMsg{Cmd: cmd, CmdArgs: cmdArgs, Text: msgText}
}

func (tc *TgChat) SetState(s *State) {
	tc.State = s
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

func (c Cmd) TgButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(c.DisplayLabel, c.Code)
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

func (tc *TgChat) TgInlineKeyboradStateCmdRows() [][]tgbotapi.InlineKeyboardButton {
	var inlineKeyboardRows [][]tgbotapi.InlineKeyboardButton

	if tc.State.WaitForWLFrgnLang {
		rowsCount := len(commons.AvailLangs)
		inlineKeyboardRows = make([][]tgbotapi.InlineKeyboardButton, rowsCount)
		for i, lang := range commons.AvailLangs {
			inlineKeyboardRows[i] = make([]tgbotapi.InlineKeyboardButton, 1)
			inlineKeyboardRows[i][0] = tgbotapi.NewInlineKeyboardButtonData(lang.Name, tc.State.StateCmd.Code+" "+lang.Code)
		}
	} else if tc.State.WaitForWLNtvLang {
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
	}

	return inlineKeyboardRows
}

func (tc *TgChat) TgInlineKeyboradMarkup() tgbotapi.InlineKeyboardMarkup {
	var inlineKeyboardMarkup tgbotapi.InlineKeyboardMarkup
	inlineKeyboardRows := make([][]tgbotapi.InlineKeyboardButton, 0)

	inlineKeyboardStateCmdRows := tc.TgInlineKeyboradStateCmdRows()
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

	return inlineKeyboardMarkup
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

func (tc *TgChat) TgOutgoingMsg() tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(tc.TgId, tc.TgOutgoingMsgText())

	// msg.ReplyToMessageID = iMsg.MessageID

	// TgChat control buttons
	msg.ReplyMarkup = tc.TgInlineKeyboradMarkup()

	return msg
}
