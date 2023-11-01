package usecases

import (
	"time"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MapToInnerTgChatAndSave(r res.Resources, outerTC *tgbotapi.Chat, u *usrPkg.User) (tc *tcPkg.Chat, err error) {
	if tc, err = r.TcRepo.TgChatByUserId(u.Id); err != nil {
		return nil, err
	}

	if tc == nil {
		state, _ := r.TcRepo.StartState()
		tc = &tcPkg.Chat{
			TgId:      int(outerTC.ID),
			UserId:    u.Id,
			State:     state,
			CreatedAt: time.Now(),
		}
		if err = r.TcRepo.SaveNewTgChat(tc); err != nil {
			return nil, err
		}
	}

	tc.User = u

	if tc.WLId != 0 {
		if tc.WL, err = r.WLRepo.WLById(tc.WLId); err != nil {
			return nil, err
		}
	}

	if tc.WordId != 0 {
		if tc.Word, err = r.WLRepo.WordById(tc.WordId); err != nil {
			return nil, err
		}
	}

	return tc, nil
}

func SetTgChatNextState(r res.Resources, tc *tcPkg.Chat, msg *tcPkg.IncMsg) (err error) {
	var nextState *tcPkg.State
	if msg.Cmd != nil {
		nextState, err = r.TcRepo.StateByCode(msg.Cmd.DestStateCode)
	} else if tc.State.NextStateCode != "" {
		nextState, err = r.TcRepo.StateByCode(tc.State.NextStateCode)
	}
	if err != nil {
		return err
	}

	tc.SetState(nextState)

	return nil
}
