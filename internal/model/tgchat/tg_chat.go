package tgchat

import (
	"errors"

	"github.com/anatoliy9697/c2vocab/internal/pkg/sliceutils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgChat struct {
	TgId      int64
	UserId    int32
	StateCode string
}

type State struct {
	Code             string
	Msg              string
	MsgHdr           string
	MsgFtr           string
	WaitForDataInput bool
	AvailCmdCodes    []string
}

type Cmd struct {
	Code      string
	DestState *State
}

func NewTgChat(tc *tgbotapi.Chat, userId int32, stateCode string) *TgChat {
	return &TgChat{TgId: tc.ID, UserId: userId, StateCode: stateCode}
}

func (s State) ValidateMessage(msg *tgbotapi.Message) error {
	if s.WaitForDataInput && msg.IsCommand() {
		return errors.New("ожидается ввод данных, не команда")
	}
	if !s.WaitForDataInput && !msg.IsCommand() {
		return errors.New("ожидается команда, не ввод данных")
	}
	if cmd := msg.CommandWithAt(); cmd != "" && !sliceutils.IsStrInSlice(s.AvailCmdCodes, cmd) {
		return errors.New("получена неожиданная команда")
	}
	return nil
}
