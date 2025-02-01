package tgchat

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"text/template"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Chat struct {
	TgId            int             `json:"tgId"`
	UserId          int             `json:"userId"`
	User            *usrPkg.User    `json:"user"`
	State           *State          `json:"state"`
	WLFrgnLang      *commons.Lang   `json:"wlFrgnLang"`
	WLNtvLang       *commons.Lang   `json:"wlNtvLang"`
	WLId            int             `json:"wlId"`
	WL              *wlPkg.WordList `json:"wl"`
	WordFrgn        string          `json:"wordFrgn"`
	WordId          int             `json:"wordId"`
	Word            *wlPkg.Word     `json:"word"`
	Words           []*wlPkg.Word   `json:"words"`
	ExcersiceCode   string          `json:"excersiceCode"`
	Excersice       *Excersice      `json:"excersice"`
	TrainedWordsIds string          `json:"trainedWordsIds"`
	PrevTaskResult  string          `json:"prevTaskReult"`
	BotLastMsgId    int             `json:"botLastMsgId"`
}

type State struct {
	Code             string             `json:"code"`
	MsgHdr           string             `json:"-"`
	MsgBody          string             `json:"-"`
	MsgFtr           string             `json:"-"`
	MsgTmpl          *template.Template `json:"-"`
	WaitForDataInput bool               `json:"-"`
	Cmd              *Cmd               `json:"-"`
	NextStateCode    string             `json:"-"`
	AvailCmds        [][]*Cmd           `json:"-"`
}

type Cmd struct {
	Code           string `json:"code"`
	DisplayLabel   string `json:"-"`
	DestStateCode  string `json:"-"`
	NotEmptyWLOnly bool   `json:"-"`
}

type Excersice struct {
	Code             string             `json:"code"`
	Name             string             `json:"-"`
	TaskText         string             `json:"-"`
	TaskTextTmpl     *template.Template `json:"-"`
	WaitForDataInput bool               `json:"-"`
	Cmd              *Cmd               `json:"-"`
}

type IncMsgValidationErr struct {
	Msg string
}

func (e IncMsgValidationErr) Error() string {
	return e.Msg
}

type IncMsg struct {
	Id            int                  `json:"id"`
	CmdCode       string               `json:"cmdCode"`
	Cmd           *Cmd                 `json:"cmd"`
	CmdArgs       []string             `json:"cmdArgs"`
	Text          string               `json:"text"`
	ValidationErr *IncMsgValidationErr `json:"-"`
}

type outMsgTmplArgs struct {
	ErrText           string
	WLName            string
	UsrTgFName        string
	UsrTgLName        string
	WLFrgnLang        string
	WLNtvLang         string
	WordsNum          int
	WordForeign       string
	WordNative        string
	ExerciseTaskText  string
	PrevTaskResult    string
	WordMemPercentage int
	WLMemPercentage   int
	WordSearchResult  string
}

var (
	ErrEmptyCmd            IncMsgValidationErr = IncMsgValidationErr{Msg: "получена пустая команда"}
	ErrUnexpectedCmd       IncMsgValidationErr = IncMsgValidationErr{Msg: "получена неожиданная команда"}
	ErrUnexpectedDataInput IncMsgValidationErr = IncMsgValidationErr{Msg: "ожидается команда, не ввод данных"}
	ErrEmptyDataInput      IncMsgValidationErr = IncMsgValidationErr{Msg: "получена пустая строка в качестве входных данных"}
)

var outMsgArgsRegExpInst *regexp.Regexp

func OutMsgArgsRegExp() *regexp.Regexp {
	if outMsgArgsRegExpInst == nil {
		outMsgArgsRegExpInst = regexp.MustCompile(`{{\.([^{]+)}}`)
	}

	return outMsgArgsRegExpInst
}

func (tc *Chat) SetState(s *State) {
	tc.State = s
}

func (tc *Chat) SetBotLastMsgId(msgId int) {
	tc.BotLastMsgId = msgId
}

func (tc *Chat) SetWords(words []*wlPkg.Word) {
	tc.Words = words
}

func (tc *Chat) WordsIdsStr() (s string) {
	for i, w := range tc.Words {
		if i == 0 {
			s += fmt.Sprintf("%d", w.Id)
		} else {
			s += fmt.Sprintf(", %d", w.Id)
		}
	}

	return
}

func (tc *Chat) OutMsgArgs(tmpl string, errText string) (args *outMsgTmplArgs, err error) {
	args = &outMsgTmplArgs{}

	submatches := OutMsgArgsRegExp().FindAllStringSubmatch(tmpl, -1)

	for _, submatch := range submatches {
		for _, group := range submatch {
			switch group {
			case "ErrText":
				if errText != "" {
					args.ErrText = errText + "\n\n"
				}
			case "UsrTgFName":
				args.UsrTgFName = tc.User.TgFirstName
			case "UsrTgLName":
				args.UsrTgLName = tc.User.TgLastName
			case "WLName":
				args.WLName = tc.WL.Name
			case "WLFrgnLang":
				args.WLFrgnLang = tc.WL.FrgnLang.Name
			case "WLNtvLang":
				args.WLNtvLang = tc.WL.NtvLang.Name
			case "WordsNum":
				args.WordsNum = tc.WL.WordsNum
			case "WordForeign":
				args.WordForeign = tc.Word.Foreign
			case "WordNative":
				args.WordNative = tc.Word.Native
			case "ExerciseTaskText":
				if args.ExerciseTaskText, err = tc.ExcersiceTaskText(); err != nil {
					return nil, err
				}
			case "PrevTaskResult":
				if tc.PrevTaskResult != "" {
					args.PrevTaskResult = tc.PrevTaskResult + "\n\n"
				}
			case "WordMemPercentage":
				args.WordMemPercentage = tc.Word.MemPercentage
			case "WLMemPercentage":
				args.WLMemPercentage = tc.WL.MemPercentage
			case "WordSearchResult":
				args.WordSearchResult = tc.WordSearchResultText()
			}
		}
	}

	return args, nil
}

func (tc *Chat) OutMsgText(errText string) (string, error) {
	var err error

	tmplText := tc.State.OutMsgTmplContent()

	var args *outMsgTmplArgs
	if args, err = tc.OutMsgArgs(tmplText, errText); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tc.State.MsgTmpl.Execute(&buf, args); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (tc *Chat) ExcersiceTaskText() (string, error) {
	var err error

	tmplText := tc.Excersice.TaskTextTmplContent()

	var args *outMsgTmplArgs
	if args, err = tc.OutMsgArgs(tmplText, ""); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = tc.Excersice.TaskTextTmpl.Execute(&buf, args); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (tc *Chat) WordSearchResultText() string {
	if len(tc.Words) == 0 {
		return "Не удалось ничего найти"
	}

	text := ""
	for i, w := range tc.Words {
		text += fmt.Sprintf("%d. %s - %s (%s)\n", (i + 1), w.Foreign, w.Native, w.WLName)
	}

	return text
}

func (tc *Chat) IsWaitForDataInput() bool {
	return tc.State.WaitForDataInput || (tc.Excersice != nil && tc.Excersice.WaitForDataInput)
}

func (tc *Chat) AddTrainedWordId(id int) {
	if tc.TrainedWordsIds == "" {
		tc.TrainedWordsIds = strconv.Itoa(id)
	} else {
		tc.TrainedWordsIds += ", " + strconv.Itoa(id)
	}
}

func (tc *Chat) SetPrevTaskResult(result string) {
	tc.PrevTaskResult = result
}

func (tc *Chat) IsCmdAvail(cmdCode string) bool {
	return tc.State.IsCmdAvail(cmdCode) || (tc.Excersice != nil && tc.Excersice.Cmd != nil && tc.Excersice.Cmd.Code == cmdCode)
}

func (s State) OutMsgTmplContent() string {
	msgTmpl := s.MsgBody
	if s.MsgHdr != "" {
		msgTmpl = s.MsgHdr + "\n\n" + msgTmpl
	}
	if s.MsgFtr != "" {
		msgTmpl += "\n\n" + s.MsgFtr
	}

	return "{{.ErrText}}" + msgTmpl
}

func (s State) IsWaitForDataInput() bool {
	return s.WaitForDataInput
}

func (s State) IsCmdAvail(cmdCode string) bool {
	for _, cmdsRow := range s.AvailCmds {
		for _, cmd := range cmdsRow {
			if cmd.Code == cmdCode {
				return true
			}
		}
	}

	if s.Cmd != nil && s.Cmd.Code == cmdCode {
		return true
	}

	return false
}

func (s State) AvailCmdsByFlgs(emptyWL bool) [][]*Cmd {
	if !emptyWL {
		return s.AvailCmds
	}

	availCmds := make([][]*Cmd, 0)

	for _, cmdsRow := range s.AvailCmds {
		tmpCmdsRow := make([]*Cmd, 0)
		for _, cmd := range cmdsRow {
			if !emptyWL || !cmd.NotEmptyWLOnly {
				tmpCmdsRow = append(tmpCmdsRow, cmd)
			}
		}
		if len(tmpCmdsRow) > 0 {
			availCmds = append(availCmds, tmpCmdsRow)
		}
	}

	return availCmds
}

func (c Cmd) TgButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(c.DisplayLabel, c.Code)
}

func (x Excersice) TaskTextTmplContent() string {
	return x.TaskText
}
