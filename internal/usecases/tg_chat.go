package usecases

import (
	"errors"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func MapToInnerTgChatAndSave(r res.Resources, outerTC *tgbotapi.Chat, u *usrPkg.User) (tc *tcPkg.Chat, err error) {
	if tc, err = r.TcRepo.TgChatByUserId(u.Id); err != nil {
		return nil, err
	}

	if tc == nil {
		state, _ := r.TcRepo.StartState()
		tc = &tcPkg.Chat{
			TgId:   int(outerTC.ID),
			UserId: u.Id,
			State:  state,
		}
		if err = r.TcRepo.SaveNewTgChat(tc); err != nil {
			return nil, err
		}
	}

	tc.User = u

	if tc.WLId != 0 {
		if tc.WL, err = r.WLRepo.WLByIdAndUserId(tc.WLId, tc.UserId); err != nil {
			return nil, err
		}
	}

	if tc.WordId != 0 {
		if tc.Word, err = r.WLRepo.WordByIdAndUserId(tc.WordId, tc.UserId); err != nil {
			return nil, err
		}
	}

	if tc.ExcersiceCode != "" {
		if tc.Excersice, err = r.TcRepo.ExcersiceByCode(tc.ExcersiceCode); err != nil {
			return nil, err
		}
	}

	return tc, nil
}

func SetTgChatNextState(r res.Resources, tc *tcPkg.Chat, msg *tcPkg.IncMsg) (err error) {
	var nextState *tcPkg.State
	if msg.Cmd != nil && msg.Cmd.DestStateCode != "" {
		nextState, err = r.TcRepo.StateByCode(msg.Cmd.DestStateCode)
	} else if tc.State.NextStateCode != "" {
		nextState, err = r.TcRepo.StateByCode(tc.State.NextStateCode)
	}
	if err != nil {
		return err
	}

	if nextState != nil {
		tc.SetState(nextState)
	}

	return nil
}

func SetTgChatExerciseNextWord(r res.Resources, tc *tcPkg.Chat) (err error) {
	if tc.Word, err = r.WLRepo.NextWordForTraining(tc.WL.Id, tc.TrainedWordsIds); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			var s *tcPkg.State
			if s, err = r.TcRepo.StateByCode("xrcs_finish"); err != nil {
				return err
			}
			tc.SetState(s)
		} else {
			return err
		}
	}
	if tc.Word != nil {
		tc.WordId = tc.Word.Id
	}

	return nil
}
