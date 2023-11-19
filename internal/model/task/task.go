package task

import (
	tcPkg "github.com/anatoliy9697/c2vocab/internal/model/tgchat"
)

type Task struct {
	Type   string      `json:"type"`
	UserId int         `json:"userId"`
	TgChat *tcPkg.Chat `json:"-"`
}
