package usecases

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func OutMsg(r res.Resources, tc *tcPkg.Chat, errText string) (msg tgbotapi.MessageConfig, err error) {
	var msgText string
	if msgText, err = tc.OutMsgText(errText); err != nil {
		return msg, err
	}

	msg = tgbotapi.NewMessage(int64(tc.TgId), msgText)

	// TgChat control buttons
	var replyMarkup tgbotapi.InlineKeyboardMarkup
	if replyMarkup, err = TgInlineKeyboradMarkup(r, tc); err != nil {
		return msg, err
	}

	msg.ReplyMarkup = replyMarkup

	return msg, nil
}

func LastMsgEditing(r res.Resources, tc *tcPkg.Chat, errText string) (editMsgConfig tgbotapi.EditMessageTextConfig, err error) {
	var msgText string
	if msgText, err = tc.OutMsgText(errText); err != nil {
		return editMsgConfig, err
	}

	// TgChat control buttons
	var replyMarkup tgbotapi.InlineKeyboardMarkup
	if replyMarkup, err = TgInlineKeyboradMarkup(r, tc); err != nil {
		return editMsgConfig, err
	}

	editMsgConfig = tgbotapi.NewEditMessageTextAndMarkup(int64(tc.TgId), tc.BotLastMsgId, msgText, replyMarkup)

	return editMsgConfig, nil
}

func SendReplyMsg(r res.Resources, tc *tcPkg.Chat, errText string) (err error) {
	var msg tgbotapi.Chattable
	if msg, err = OutMsg(r, tc, errText); err != nil {
		return err
	}

	r.Logger.Info("Sending reply message", "replyMsg", msg)

	var msgInTg tgbotapi.Message
	if msgInTg, err = r.TgBotAPI.Send(msg); err != nil {
		return err
	}

	if tc.BotLastMsgId != 0 {
		DeleteMsgInTg(r, tc.TgId, tc.BotLastMsgId)
	}

	tc.SetBotLastMsgId(msgInTg.MessageID)

	return nil
}

func EditLastMsg(r res.Resources, tc *tcPkg.Chat, errText string) (err error) {
	var msg tgbotapi.Chattable
	if msg, err = LastMsgEditing(r, tc, errText); err != nil {
		return err
	}

	r.Logger.Info("Sending reply message", "replyMsg", msg)

	if _, err = r.TgBotAPI.Send(msg); err != nil {
		return err
	}

	return nil
}
