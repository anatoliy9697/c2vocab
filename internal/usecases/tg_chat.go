package usecases

import (
	"time"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	tcRepo "github.com/anatoliy9697/c2vocab/internal/model/tgchat/repo"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	wlRepo "github.com/anatoliy9697/c2vocab/internal/model/wordlist/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MapToInnerTgChatAndSave(tcR tcRepo.Repo, wlR wlRepo.Repo, outerTC *tgbotapi.Chat, u *usrPkg.User) (tc *tcPkg.Chat, err error) {
	if tc, err = tcR.TgChatByUserId(u.Id); err != nil {
		return nil, err
	}

	if tc == nil {
		state, _ := tcR.StartState()
		tc = &tcPkg.Chat{
			TgId:      outerTC.ID,
			UserId:    u.Id,
			State:     state,
			CreatedAt: time.Now(),
		}
		if err = tcR.SaveNewTgChat(tc); err != nil {
			return nil, err
		}
	}

	if tc.WLId != 0 {
		if tc.WL, err = wlR.WLById(tc.WLId); err != nil {
			return nil, err
		}
	}

	return tc, nil
}

func SetTgChatNextState(tcR tcRepo.Repo, tc *tcPkg.Chat, msg *tcPkg.IncMsg) (err error) {
	var nextState *tcPkg.State
	if msg.Cmd != nil {
		nextState, err = tcR.StateByCode(msg.Cmd.DestStateCode)
	} else if tc.State.NextStateCode != "" {
		nextState, err = tcR.StateByCode(tc.State.NextStateCode)
	}
	if err != nil {
		return err
	}

	tc.SetState(nextState)

	return nil
}
