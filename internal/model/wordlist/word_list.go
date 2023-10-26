package wordlist

import (
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
)

type WordList struct {
	Id        int           `json:"id"`
	Active    bool          `json:"active"`
	Name      string        `json:"name"`
	FrgnLang  *commons.Lang `json:"frgnLang"`
	NtvLang   *commons.Lang `json:"ntvLang"`
	WordsNum  int           `json:"wordsNum"`
	OwnerId   int           `json:"ownerId"`
	CreatedAt time.Time     `json:"createdAt"`
}

type Word struct {
	Id        int       `json:"id"`
	Foreign   string    `json:"foreign"`
	Native    string    `json:"native"`
	WLId      int       `json:"wlId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (wl *WordList) Deactivate() {
	wl.Active = false
}
