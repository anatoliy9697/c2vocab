package wordlist

import (
	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	"github.com/anatoliy9697/c2vocab/internal/model/user"
)

type WordList struct {
	Id          uint32
	Code        string
	Name        string
	NativeLang  *commons.Lang
	ForeignLang *commons.Lang
	Words       []*Word
	Owner       *user.User
}

type Word struct {
	Id      uint32
	Native  string
	Foreign string
}
