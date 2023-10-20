package usecases

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	wlRepo "github.com/anatoliy9697/c2vocab/internal/model/wordlist/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TgOutMsg(wlR wlRepo.Repo, tc *tcPkg.Chat) (msg tgbotapi.MessageConfig, err error) {
	var msgText string
	if msgText, err = tc.OutMsgText(); err != nil {
		return msg, err
	}

	msg = tgbotapi.NewMessage(tc.TgId, msgText)

	// TgChat control buttons
	var replyMarkup tgbotapi.InlineKeyboardMarkup
	if replyMarkup, err = TgInlineKeyboradMarkup(wlR, tc); err != nil {
		return msg, err
	}

	msg.ReplyMarkup = replyMarkup

	return msg, nil
}

func TgMsgEditing(wlR wlRepo.Repo, tc *tcPkg.Chat) (editMsgConfig tgbotapi.EditMessageTextConfig, err error) {
	var msgText string
	if msgText, err = tc.OutMsgText(); err != nil {
		return editMsgConfig, err
	}

	// TgChat control buttons
	var replyMarkup tgbotapi.InlineKeyboardMarkup
	if replyMarkup, err = TgInlineKeyboradMarkup(wlR, tc); err != nil {
		return editMsgConfig, err
	}

	editMsgConfig = tgbotapi.NewEditMessageTextAndMarkup(tc.TgId, tc.BotLastMsgId, msgText, replyMarkup)

	return editMsgConfig, nil
}

func SendReplyMessage(wlR wlRepo.Repo, tgClient *tgbotapi.BotAPI, tc *tcPkg.Chat) (err error) {
	var msg tgbotapi.Chattable
	if tc.BotLastMsgId != 0 {
		msg, err = TgMsgEditing(wlR, tc)
	} else {
		msg, err = TgOutMsg(wlR, tc)
	}
	if err != nil {
		return err
	}

	var msgInTg tgbotapi.Message
	if msgInTg, err = tgClient.Send(msg); err != nil && tc.BotLastMsgId != 0 {
		if msg, err = TgOutMsg(wlR, tc); err != nil {
			return err
		}
		if msgInTg, err = tgClient.Send(msg); err != nil {
			return err
		}
	}

	tc.SetBotLastMsgId(msgInTg.MessageID)

	return nil
}
