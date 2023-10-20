package tgchat

import (
	"bytes"
	"errors"
	"regexp"
	"text/template"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Chat struct {
	TgId         int64
	UserId       int32
	CreatedAt    time.Time
	State        *State
	WLFrgnLang   *commons.Lang
	WLNtvLang    *commons.Lang
	WLId         int32
	WL           *wlPkg.WordList
	BotLastMsgId int
}

type State struct {
	Code              string
	MsgHdr            string
	MsgBody           string
	MsgFtr            string
	MsgTmpl           *template.Template
	WaitForWLFrgnLang bool
	WaitForWLNtvLang  bool
	WaitForWLName     bool
	StateCmd          *Cmd
	NextStateCode     string
	AvailCmds         [][]*Cmd
}

type Cmd struct {
	Code          string
	DisplayLabel  string
	DestStateCode string
}

type IncMsg struct {
	Id      int
	Cmd     *Cmd
	CmdArgs []string
	Text    string
}

type outMsgTmplArgs struct {
	WLName string
}

var (
	ErrEmptyCmd            = errors.New("получена пустая команда")
	ErrUnexpectedCmd       = errors.New("получена неожиданная команда")
	ErrUnexpectedDataInput = errors.New("ожидается команда, не ввод данных")
	ErrEmptyDataInput      = errors.New("получена пустая строка в качестве входных данных")
)

var outMsgArgsRegExpInst *regexp.Regexp

func OutMsgArgsRegExp() *regexp.Regexp {
	if outMsgArgsRegExpInst == nil {
		outMsgArgsRegExpInst = regexp.MustCompile(`{{\.(.*)}}`)
	}

	return outMsgArgsRegExpInst
}

func (tc *Chat) SetState(s *State) {
	tc.State = s
}

func (tc *Chat) SetBotLastMsgId(msgId int) {
	tc.BotLastMsgId = msgId
}

func (tc *Chat) OutMsgArgs(tmpl string) *outMsgTmplArgs {
	args := &outMsgTmplArgs{}

	submatches := OutMsgArgsRegExp().FindAllStringSubmatch(tmpl, -1)

	for _, submatch := range submatches {
		for _, group := range submatch {
			switch group {
			case "WLName":
				args.WLName = tc.WL.Name
			}
		}
	}

	return args
}

func (tc *Chat) OutMsgText() (string, error) {
	var err error

	tmplText := tc.State.OutMsgTmplContent()

	var buf bytes.Buffer
	if err = tc.State.MsgTmpl.Execute(&buf, tc.OutMsgArgs(tmplText)); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s State) OutMsgTmplContent() string {
	msgTmpl := s.MsgBody
	if s.MsgHdr != "" {
		msgTmpl = s.MsgHdr + "\n\n" + msgTmpl
	}
	if s.MsgFtr != "" {
		msgTmpl += "\n\n" + s.MsgFtr
	}

	return msgTmpl
}

func (s State) IsWaitForDataInput() bool {
	return s.WaitForWLName
}

func (s State) IsCmdAvail(cmdCode string) bool {
	for _, cmdsRow := range s.AvailCmds {
		for _, cmd := range cmdsRow {
			if cmd.Code == cmdCode {
				return true
			}
		}
	}

	if s.StateCmd != nil && s.StateCmd.Code == cmdCode {
		return true
	}

	return false
}

func (s State) TgInlineKeyboradAvailCmdsRows() [][]tgbotapi.InlineKeyboardButton {
	inlineKeyboardRows := make([][]tgbotapi.InlineKeyboardButton, len(s.AvailCmds))
	for i, cmdsRow := range s.AvailCmds {
		inlineKeyboardRows[i] = make([]tgbotapi.InlineKeyboardButton, len(cmdsRow))
		for j, cmd := range cmdsRow {
			inlineKeyboardRows[i][j] = cmd.TgButton()
		}
	}

	return inlineKeyboardRows
}

func (c Cmd) TgButton() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(c.DisplayLabel, c.Code)
}
