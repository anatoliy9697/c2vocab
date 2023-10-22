package commons

type Lang struct {
	Code string `json:"code"`
	Name string `json:"-"`
}

var AvailLangs = []*Lang{
	{"en", "English"},
	{"ru", "Русский"},
}

func LangByCode(code string) *Lang {
	for _, lang := range AvailLangs {
		if lang.Code == code {
			return lang
		}
	}

	return nil
}
