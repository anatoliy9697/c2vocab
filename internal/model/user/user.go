package user

import (
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	Id          int32
	TgId        int64
	TgUserName  string
	TgFirstName string
	TgLastName  string
	Lang        *commons.Lang
	TgIsBot     bool
	CreatedAt   time.Time
}

func MapToInner(u *tgbotapi.User) *User {
	return &User{
		TgId:        u.ID,
		TgUserName:  u.UserName,
		TgFirstName: u.FirstName,
		TgLastName:  u.LastName,
		Lang:        commons.LangByCode(u.LanguageCode),
		TgIsBot:     u.IsBot,
		CreatedAt:   time.Now(),
	}
}
