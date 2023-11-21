package usecases

import (
	tskPkg "github.com/anatoliy9697/c2vocab/internal/model/task"
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
)

func HandleTask(r res.Resources, task tskPkg.Task) (err error) {
	switch task.Type {
	case "to_main_menu":
		err = HandleReturnToMainMenu(r, task)
	}

	return err
}

func HandleReturnToMainMenu(r res.Resources, task tskPkg.Task) (err error) {
	var tc *tcPkg.Chat
	if tc, err = r.TcRepo.TgChatByUserId(task.UserId); err != nil {
		return err
	}

	if tc.User, err = r.UsrRepo.ById(tc.UserId); err != nil {
		return nil
	}

	ClearTgChaTmpFields(tc)

	var s *tcPkg.State
	if s, err = r.TcRepo.StateByCode("main_menu"); err != nil {
		return err
	}
	tc.SetState(s)

	if tc.BotLastMsgId != 0 {
		err = EditLastMsg(r, tc, "")
	} else {
		err = SendReplyMsg(r, tc, "")
	}
	if err != nil {
		return err
	}

	err = r.TcRepo.UpdateTgChat(tc, false)

	return err
}
