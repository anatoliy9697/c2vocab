package resources

import (
	"log/slog"

	tsk "github.com/anatoliy9697/c2vocab/internal/model/task/repo"
	tgchat "github.com/anatoliy9697/c2vocab/internal/model/tgchat/repo"
	usr "github.com/anatoliy9697/c2vocab/internal/model/user/repo"
	wl "github.com/anatoliy9697/c2vocab/internal/model/wordlist/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type LockConfig struct {
	AttemptsAmount     int
	TimeForNextAttempt int // ms
	TimeForReassign    int // ms
}

type Resources struct {
	UsrRepo  usr.Repo
	TcRepo   tgchat.Repo
	WLRepo   wl.Repo
	TskRepo  tsk.Repo
	TgBotAPI *tgbotapi.BotAPI
	Logger   *slog.Logger
	LockConf LockConfig
}
