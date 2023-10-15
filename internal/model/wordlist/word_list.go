package wordlist

import (
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
)

type WordList struct {
	Id     int32
	Active bool
	// Code        string
	Name      string
	FrgnLang  *commons.Lang
	NtvLang   *commons.Lang
	OwnerId   int32
	CreatedAt time.Time
}

type Word struct {
	Id      int32
	Native  string
	Foreign string
}
