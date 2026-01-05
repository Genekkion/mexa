package chatdomain

type ReplyMarkup struct {
	InlineKeyboard [][]InlineKeyboardEntry `json:"inline_keyboard"`
}

type InlineKeyboardEntry struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}
type SendConfig struct {
	Rmk       *ReplyMarkup
	ParseMode string
}

type SendOpt func(*SendConfig)

func WithReplyMarkup(rmk ReplyMarkup) SendOpt {
	return func(sc *SendConfig) {
		sc.Rmk = &rmk
	}
}

func WithParseMode(pm string) SendOpt {
	return func(sc *SendConfig) {
		sc.ParseMode = pm
	}
}
