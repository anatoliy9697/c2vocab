package tgchat

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgChat struct {
	TgId      int64
	UserId    int32
	StateCode string
}

type State struct {
	Code          string
	Msg           string
	MsgHdr        string
	MsgFtr        string
	AvailCmdCodes []string
}

type Cmd struct {
	Code      string
	DestState *State
}

func MapToInner(tc *tgbotapi.Chat) *TgChat {
	return &TgChat{TgId: tc.ID}
}
