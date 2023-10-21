package resources

import (
	"log/slog"

	tgchat "github.com/anatoliy9697/c2vocab/internal/model/tgchat/repo"
	usr "github.com/anatoliy9697/c2vocab/internal/model/user/repo"
	wl "github.com/anatoliy9697/c2vocab/internal/model/wordlist/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Resources struct {
	UsrRepo  usr.Repo
	TcRepo   tgchat.Repo
	WLRepo   wl.Repo
	TgBotAPI *tgbotapi.BotAPI
	Logger   *slog.Logger
}
