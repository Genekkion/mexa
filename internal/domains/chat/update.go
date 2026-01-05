package chatdomain

type Update struct {
	UpdateId      int `json:"update_id"`
	CallbackQuery *struct {
		Id      string `json:"id"`
		From    User   `json:"from"`
		Message struct {
			Message

			Date        int         `json:"date"`
			ReplyMarkup ReplyMarkup `json:"reply_markup"`
		} `json:"message"`
		ChatInstance string `json:"chat_instance"`
		Data         string `json:"data"`
	} `json:"callback_query"`
	Message *Message `json:"message"`
}

func (u Update) ToIgnore() bool {
	return u.CallbackQuery == nil &&
		u.Message == nil
}

type User struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Message struct {
	MessageId int  `json:"message_id"`
	From      User `json:"from"`
	Chat      struct {
		Id        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Date int    `json:"date"`
	Text string `json:"text"`
}

func (u Update) ChatId() ChatId {
	if u.IsCallback() {
		return u.CallbackQuery.Message.Chat.Id
	}
	return u.Message.Chat.Id
}

func (u Update) UserId() UserId {
	if u.IsCallback() {
		return u.CallbackQuery.From.Id
	}
	return u.Message.From.Id
}

func (u Update) Username() string {
	if u.IsCallback() {
		return u.CallbackQuery.From.Username
	}
	return u.Message.From.Username
}

func (u Update) IsCommand() bool {
	if !u.IsMessage() {
		return false
	} else if len(u.Message.Text) == 0 {
		return false
	}
	return u.Message.Text[0] == '/'
}

func (u Update) IsCallback() bool {
	return u.CallbackQuery != nil
}

func (u Update) IsMessage() bool {
	return u.Message != nil
}
