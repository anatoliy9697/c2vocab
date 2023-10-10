package tgchat

import (
	"errors"

	"github.com/anatoliy9697/c2vocab/pkg/sliceutils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgChat struct {
	TgId   int64
	UserId int32
	State  *State
}

type State struct {
	Code   string
	MsgHdr string
	Msg    string
	MsgFtr string
	// WaitForDataInput bool
	AvailCmdCodes []string
}

type Cmd struct {
	Code      string
	DestState *State
}

func NewTgChat(tc *tgbotapi.Chat, userId int32, state *State) *TgChat {
	return &TgChat{
		TgId:   tc.ID,
		UserId: userId,
		State:  state,
	}
}

func (tc *TgChat) ValidateMessage(msg *tgbotapi.Message) error {
	// if tc.State.WaitForDataInput && msg.IsCommand() {
	// 	return errors.New("ожидается ввод данных, не команда")
	// }
	// if !tc.State.WaitForDataInput && !msg.IsCommand() {
	// 	return errors.New("ожидается команда, не ввод данных")
	// }
	if cmd := msg.CommandWithAt(); cmd == "" || !sliceutils.IsStrInSlice(tc.State.AvailCmdCodes, cmd) {
		return errors.New("получена неожиданная команда")
	}
	return nil
}

func (tc *TgChat) SetState(s *State) {
	tc.State = s
}
