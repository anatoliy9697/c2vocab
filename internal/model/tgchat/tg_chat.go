package tgchat

import (
	"errors"
)

type TgChat struct {
	TgId   int64
	UserId int32
	State  *State
}

type State struct {
	Code             string
	MsgHdr           string
	Msg              string
	MsgFtr           string
	WaitForDataInput bool
	AvailCmdCodes    []string
}

type Cmd struct {
	Code      string
	DestState *State
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

func (s *State) IsWaitForDataInput() bool {
	return s.WaitForDataInput
}
