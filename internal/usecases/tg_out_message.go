package usecases

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TgOutMsg(r res.Resources, tc *tcPkg.Chat, errText string) (msg tgbotapi.MessageConfig, err error) {
	var msgText string
	if msgText, err = tc.OutMsgText(errText); err != nil {
		return msg, err
	}

	msg = tgbotapi.NewMessage(tc.TgId, msgText)

	// TgChat control buttons
	var replyMarkup tgbotapi.InlineKeyboardMarkup
	if replyMarkup, err = TgInlineKeyboradMarkup(r, tc); err != nil {
		return msg, err
	}

	msg.ReplyMarkup = replyMarkup

	return msg, nil
}

func TgMsgEditing(r res.Resources, tc *tcPkg.Chat, errText string) (editMsgConfig tgbotapi.EditMessageTextConfig, err error) {
	var msgText string
	if msgText, err = tc.OutMsgText(errText); err != nil {
		return editMsgConfig, err
	}

	// TgChat control buttons
	var replyMarkup tgbotapi.InlineKeyboardMarkup
	if replyMarkup, err = TgInlineKeyboradMarkup(r, tc); err != nil {
		return editMsgConfig, err
	}

	editMsgConfig = tgbotapi.NewEditMessageTextAndMarkup(tc.TgId, tc.BotLastMsgId, msgText, replyMarkup)

	return editMsgConfig, nil
}

func SendReplyMsg(r res.Resources, tc *tcPkg.Chat, errText string) (err error) {
	var msg tgbotapi.Chattable
	if msg, err = TgOutMsg(r, tc, errText); err != nil {
		return err
	}

	r.Logger.Info("Sending reply message", "replyMsg", msg)

	var msgInTg tgbotapi.Message
	if msgInTg, err = r.TgBotAPI.Send(msg); err != nil && tc.BotLastMsgId != 0 {
		if msg, err = TgOutMsg(r, tc, errText); err != nil {
			return err
		}
		if msgInTg, err = r.TgBotAPI.Send(msg); err != nil {
			return err
		}
	}

	if tc.BotLastMsgId != 0 {
		DeleteMsgInTg(r, tc.TgId, tc.BotLastMsgId)
	}

	tc.SetBotLastMsgId(msgInTg.MessageID)

	return nil
}
