package usecases

import (
	"strings"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ValidateIncMsgAndMapToInner(r res.Resources, tc *tcPkg.Chat, upd *tgbotapi.Update) (incMsg *tcPkg.IncMsg, err error) {
	var id int
	var cmdCode, msgText string
	var cmdArgs []string

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
		id = upd.Message.MessageID
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

	var cmd *tcPkg.Cmd
	if cmdCode != "" {
		if cmd, err = r.TcRepo.CmdByCode(cmdCode); err != nil {
			return nil, err
		}
	}

	return &tcPkg.IncMsg{
		Id:      id,
		Cmd:     cmd,
		CmdArgs: cmdArgs,
		Text:    msgText,
	}, nil
}

func DeleteMsgInTg(r res.Resources, chatId int64, msgId int) (err error) {
	delMsg := tgbotapi.NewDeleteMessage(chatId, msgId)

	_, err = r.TgBotAPI.Send(delMsg)

	return err
}

func ProcessIncMsg(r res.Resources, tc *tcPkg.Chat, msg *tcPkg.IncMsg) (err error) {
	switch {
	case msg.Cmd != nil && (msg.Cmd.Code == "start" || msg.Cmd.Code == "to_main_menu"):
		ClearTgChaTmpFields(tc)
	case msg.Cmd != nil && msg.Cmd.Code == "wl_frgn_lang":
		SetTgChatWLFrgnLang(tc, msg.CmdArgs[0])
	case msg.Cmd != nil && msg.Cmd.Code == "wl_ntv_lang":
		SetTgChatWLNtvLang(tc, msg.CmdArgs[0])
	case msg.Cmd != nil && msg.Cmd.Code == "wl":
		if err = SetTgChatWL(r, tc, msg.CmdArgs[0]); err != nil {
			return err
		}
	case msg.Cmd != nil && msg.Cmd.Code == "delete_wl":
		if err = DeleteWL(r, tc.WL); err != nil {
			return err
		}
	case msg.Text != "" && tc.State.WaitForWLName:
		if err = CreateWL(r, tc, msg.Text); err != nil {
			return err
		}
	}

	return nil
}

func ProcessUpd(r res.Resources, tc *tcPkg.Chat, upd *tgbotapi.Update) (err error) {
	var msg *tcPkg.IncMsg
	if msg, err = ValidateIncMsgAndMapToInner(r, tc, upd); err != nil {
		return err
	}

	// Removing incoming msg in tgChat, if it's not callback query
	if msg.Id != 0 {
		DeleteMsgInTg(r, tc.TgId, msg.Id)
	}

	if err = ProcessIncMsg(r, tc, msg); err != nil {
		return err
	}

	if err = SetTgChatNextState(r, tc, msg); err != nil {
		return err
	}

	return nil
}
