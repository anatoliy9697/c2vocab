package wordlist

import (
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
)

type WordList struct {
	Id        int32         `json:"id"`
	Active    bool          `json:"active"`
	Name      string        `json:"name"`
	FrgnLang  *commons.Lang `json:"frgnLang"`
	NtvLang   *commons.Lang `json:"ntvLang"`
	OwnerId   int32         `json:"ownerId"`
	CreatedAt time.Time     `json:"createdAt"`
}

type Word struct {
	Id      int32
	Native  string
	Foreign string
}

func (wl *WordList) Deactivate() {
	wl.Active = false
}
