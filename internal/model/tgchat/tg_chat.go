package tgchat

import (
	"bytes"
	"regexp"
	"text/template"
	"time"

	"github.com/anatoliy9697/c2vocab/internal/model/commons"
	usrPkg "github.com/anatoliy9697/c2vocab/internal/model/user"
	wlPkg "github.com/anatoliy9697/c2vocab/internal/model/wordlist"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Chat struct {
	TgId         int             `json:"tgId"`
	UserId       int             `json:"userId"`
	User         *usrPkg.User    `json:"user"`
	CreatedAt    time.Time       `json:"createdAt"`
	State        *State          `json:"-"`
	WLFrgnLang   *commons.Lang   `json:"wlFrgnLang"`
	WLNtvLang    *commons.Lang   `json:"wlNtvLang"`
	WLId         int             `json:"wlId"`
	WL           *wlPkg.WordList `json:"wl"`
	WordFrgn     string          `json:"wordFrgn"`
	WordId       int             `json:"wordId"`
	Word         *wlPkg.Word     `json:"word"`
	BotLastMsgId int             `json:"botLastMsgId"`
}

type State struct {
	Code             string             `json:"code"`
	MsgHdr           string             `json:"-"`
	MsgBody          string             `json:"-"`
	MsgFtr           string             `json:"-"`
	MsgTmpl          *template.Template `json:"-"`
	WaitForDataInput bool               `json:"-"`
	StateCmd         *Cmd               `json:"-"`
	NextStateCode    string             `json:"-"`
	AvailCmds        [][]*Cmd           `json:"-"`
}

type Cmd struct {
	Code          string `json:"code"`
	DisplayLabel  string `json:"-"`
	DestStateCode string `json:"-"`
}

type Excersice struct {
	Code string `json:"code"`
	Name string `json:"name"`
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
	ErrText     string
	WLName      string
	UsrTgFName  string
	UsrTgLName  string
	WLFrgnLang  string
	WLNtvLang   string
	WordsNum    int
	WordForeign string
	WordNative  string
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

func (tc *Chat) OutMsgArgs(tmpl string, errText string) *outMsgTmplArgs {
	args := &outMsgTmplArgs{}

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
			}
		}
	}

	return args
}

func (tc *Chat) OutMsgText(errText string) (string, error) {
	var err error

	tmplText := tc.State.OutMsgTmplContent()

	var buf bytes.Buffer
	if err = tc.State.MsgTmpl.Execute(&buf, tc.OutMsgArgs(tmplText, errText)); err != nil {
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

	if s.StateCmd != nil && s.StateCmd.Code == cmdCode {
		return true
	}

	return false
}

func (s State) TgInlineKeyboradAvailCmdsRows() [][]tgbotapi.InlineKeyboardButton { // TODO: Не отображать некоторые недоступные кнопки
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
