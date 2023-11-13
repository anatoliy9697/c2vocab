package control

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	"github.com/anatoliy9697/c2vocab/internal/usecases"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type EventHandler struct {
	Code string
	Res  res.Resources
}

func (eh EventHandler) Run(done chan string, upd *tgbotapi.Update) {
	defer func() { done <- eh.Code }()

	var err error
	var usr *usrPkg.User
	var tc *tcPkg.Chat

	eh.Res.Logger = eh.Res.Logger.With("handlerCode", eh.Code)

	defer func() {
		usecases.ProcessErr(eh.Res, tc, usr, err)
	}()

	if usr, err = usecases.MapToInnerUserAndSave(eh.Res, upd.SentFrom()); err != nil {
		return
	}

	eh.Res.Logger.Debug("User mapped to inner", "user", usr)

	if tc, err = usecases.MapToInnerTgChatAndSave(eh.Res, upd.FromChat(), usr); err != nil {
		return
	}

	eh.Res.Logger.Debug("TgChat mapped to inner", "tgChat", tc)

	// Ignore non-message and non-command events
	if upd.Message == nil && upd.CallbackQuery == nil {
		return
	}

	if err = usecases.ProcessUpd(eh.Res, tc, upd); err != nil {
		return
	}

	eh.Res.Logger.Debug("Update processed", "tgChat", tc)

	if err = usecases.SendReplyMsg(eh.Res, tc, ""); err != nil {
		return
	}

	eh.Res.Logger.Info("TgChat updating in DB", "tgChat", tc)

	if err = eh.Res.TcRepo.UpdateTgChat(tc, true); err != nil {
		return
	}
}
