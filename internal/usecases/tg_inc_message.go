package usecases

import (
	"strings"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MapIncMsgToInner(r res.Resources, tc *tcPkg.Chat, upd *tgbotapi.Update) (msg *tcPkg.IncMsg) {
	msg = &tcPkg.IncMsg{}

	if upd.CallbackQuery != nil {
		if upd.CallbackQuery.Data == "" {
			msg.ValidationErr = &tcPkg.ErrEmptyCmd
			return msg
		}
		cmdParts := strings.Split(upd.CallbackQuery.Data, " ")
		if len(cmdParts) > 1 {
			msg.CmdArgs = cmdParts[1:]
		}
		msg.CmdCode = strings.ReplaceAll(cmdParts[0], "/", "")
		if !tc.IsCmdAvail(msg.CmdCode) {
			msg.ValidationErr = &tcPkg.ErrUnexpectedCmd
			return msg
		}
	} else if upd.Message != nil {
		msg.Id = upd.Message.MessageID
		if upd.Message.IsCommand() {
			msg.CmdCode = upd.Message.Command()
			msg.CmdArgs = strings.Split(upd.Message.CommandArguments(), " ")
			if msg.CmdCode != "start" {
				msg.ValidationErr = &tcPkg.ErrUnexpectedCmd
				return msg
			}
		} else {
			msg.Text = upd.Message.Text
			if !tc.IsWaitForDataInput() {
				msg.ValidationErr = &tcPkg.ErrUnexpectedDataInput
				return msg
			}
		}
	}

	if msg.CmdCode != "" {
		var err error
		msg.Cmd, err = r.TcRepo.CmdByCode(msg.CmdCode)
		if err != nil {
			msg.ValidationErr = &tcPkg.IncMsgValidationErr{Msg: err.Error()}
		}
	}

	return msg
}

func DeleteMsgInTg(r res.Resources, chatId int, msgId int) (err error) {
	delMsg := tgbotapi.NewDeleteMessage(int64(chatId), msgId)

	_, err = r.TgBotAPI.Send(delMsg)

	return err
}

func ProcessIncMsg(r res.Resources, tc *tcPkg.Chat, msg *tcPkg.IncMsg) (err error) {
	switch {

	// Navigation
	case msg.Cmd != nil && (msg.Cmd.Code == "start" || msg.Cmd.Code == "to_main_menu"):
		ClearTgChaTmpFields(tc)
	case msg.Cmd != nil && msg.Cmd.Code == "back_to_all_wl":
		ClearWLFields(tc)
	case msg.Cmd != nil && msg.Cmd.Code == "back_to_wl":
		ClearWordCreationFields(tc)
	case msg.Cmd != nil && msg.Cmd.Code == "back_to_all_w":
		ClearWordFields(tc)
	case msg.Cmd != nil && msg.Cmd.Code == "finish_xrcs":
		ClearExerciseFields(tc)

	// Word list
	case msg.Cmd != nil && msg.Cmd.Code == "wl":
		if err = SetTgChatWL(r, tc, msg.CmdArgs[0]); err != nil {
			return err
		}
	case msg.Cmd != nil && (msg.Cmd.Code == "wl_creation_frgn_lang" || msg.Cmd.Code == "wl_editing_frgn_lang"):
		SetTgChatWLFrgnLang(tc, msg.CmdArgs[0])
	case msg.Cmd != nil && (msg.Cmd.Code == "wl_creation_ntv_lang" || msg.Cmd.Code == "wl_editing_ntv_lang"):
		SetTgChatWLNtvLang(tc, msg.CmdArgs[0])
	case msg.Text != "" && tc.State.Code == "wl_creation_name":
		if err = CreateWL(r, tc, msg.Text); err != nil {
			return err
		}
	case msg.Text != "" && tc.State.Code == "wl_editing_name":
		if err = EditWL(r, tc, msg.Text); err != nil {
			return err
		}
	case msg.Cmd != nil && msg.Cmd.Code == "confirm_wl_del":
		if err = DeleteWL(r, tc.WL); err != nil {
			return err
		}

	// Word control
	case msg.Cmd != nil && msg.Cmd.Code == "w":
		if err = SetTgChatWord(r, tc, msg.CmdArgs[0]); err != nil {
			return err
		}
	case msg.Text != "" && tc.State.Code == "w_addition_frgn":
		SetTgChatWordFrgn(tc, msg.Text)
	case msg.Text != "" && tc.State.Code == "w_addition_ntv":
		if err = CreateWord(r, tc, msg.Text); err != nil {
			return err
		}
	case msg.Cmd != nil && msg.Cmd.Code == "confirm_w_del":
		if err = DeleteWord(r, tc.Word); err != nil {
			return err
		}

	// Learning
	case msg.Cmd != nil && msg.Cmd.Code == "xrcs":
		if err = SetTgChatExercise(r, tc, msg.CmdArgs[0]); err != nil {
			return err
		}
	case msg.Text != "" && tc.State.Code == "xrcs":
		if err = ProcessUserTaskDataInput(r, tc, msg.Text); err != nil {
			return err
		}
	case msg.Cmd != nil && msg.Cmd.Code == "ans" && tc.State.Code == "xrcs":
		if err = ProcessUserTaskAnswer(r, tc, msg.CmdArgs[0]); err != nil {
			return err
		}

	default:

	}

	return nil
}

func ProcessUpd(r res.Resources, tc *tcPkg.Chat, upd *tgbotapi.Update) (err error) {
	msg := MapIncMsgToInner(r, tc, upd)

	r.Logger.Info("Got incoming message", "incMsg", msg)

	// Removing incoming msg in tgChat, if it's not callback query
	if msg.Id != 0 {
		DeleteMsgInTg(r, tc.TgId, msg.Id)
	}

	if msg.ValidationErr != nil {
		return msg.ValidationErr
	}

	if err = ProcessIncMsg(r, tc, msg); err != nil {
		return err
	}

	if err = SetTgChatNextState(r, tc, msg); err != nil {
		return err
	}

	return nil
}
