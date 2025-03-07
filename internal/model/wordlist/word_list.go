package wordlist

import (
	"math/rand"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

type WordList struct {
	Id            int           `json:"id"`
	Active        bool          `json:"active"`
	Name          string        `json:"name"`
	FrgnLang      *commons.Lang `json:"frgnLang"`
	NtvLang       *commons.Lang `json:"ntvLang"`
	WordsNum      int           `json:"wordsNum"`
	MemPercentage int           `json:"memPercentage"`
	OwnerId       int           `json:"ownerId"`
	CreatedAt     time.Time     `json:"createdAt"`
}

type Word struct {
	Id            int       `json:"id"`
	Active        bool      `json:"active"`
	Foreign       string    `json:"foreign"`
	Native        string    `json:"native"`
	MemPercentage int       `json:"memPercentage"`
	WLId          int       `json:"wlId"`
	WLName        string    `json:"wlName"`
	CreatedAt     time.Time `json:"createdAt"`
}

type AnswerOption struct {
	Answer    string
	IsCorrect string // "1" - correct, "0" - incorrect
}

func MixAnswerOptions(opts []AnswerOption) []AnswerOption {
	for i := len(opts) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		opts[i], opts[j] = opts[j], opts[i]
	}

	return opts
}

func (wl *WordList) Deactivate() {
	wl.Active = false
}

func (w *Word) Deactivate() {
	w.Active = false
}

func (w *Word) MemPercentageDowngrade() int {
	if w.MemPercentage == 0 {
		return -1
	}

	memPercentage := w.MemPercentage
	dg := 0
	for ; memPercentage < 50; memPercentage *= 2 {
		dg++
	}

	return dg
}
