package usecases

import (
	"errors"

	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func MapToInnerTgChatAndSaveWithLocking(r res.Resources, outerTC *tgbotapi.Chat, u *usrPkg.User, handlerCode string) (tc *tcPkg.Chat, err error) {
	var chatExists bool
	if chatExists, err = r.TcRepo.IsChatExistsByUserId(u.Id); err != nil {
		return nil, err
	}

	if chatExists {
		if err = r.TcRepo.LockChatByUserId(u.Id, handlerCode, r.LockConf.TimeForReassign, r.LockConf.AttemptsAmount, r.LockConf.TimeForNextAttempt); err != nil {
			return nil, err
		}
		tc, err = r.TcRepo.ChatByUserId(u.Id)
	} else {
		state, _ := r.TcRepo.StartState()
		tc = &tcPkg.Chat{
			TgId:   int(outerTC.ID),
			UserId: u.Id,
			State:  state,
		}
		err = r.TcRepo.SaveNewChat(tc, handlerCode)
	}
	if err != nil {
		return nil, err
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

func UnlockTgChat(r res.Resources, u *usrPkg.User) {
	if u != nil {
		if err := r.TcRepo.UnlockChatByUserId(u.Id); err != nil {
			r.Logger.Error(err.Error())
		}
		r.Logger.Debug("TgChat unlocked")
	}
}
