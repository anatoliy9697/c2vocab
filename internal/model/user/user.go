package user

import (
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	Id          int32         `json:"id"`
	TgId        int64         `json:"tgId"`
	TgUserName  string        `json:"tgUserName"`
	TgFirstName string        `json:"tgFistName"`
	TgLastName  string        `json:"tgLastName"`
	Lang        *commons.Lang `json:"lang"`
	TgIsBot     bool          `json:"tgIsBot"`
	CreatedAt   time.Time     `json:"createdAt"`
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
