package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	chatdomain "mexa/internal/domains/chat"
	mexadomain "mexa/internal/domains/mexa"
	"mexa/internal/utils/set"
	"net/http"
	"sync"
)

const (
	telegramBaseUrl = "https://api.telegram.org/bot"
)

type Bot struct {
	token   string
	baseUrl string

	client *http.Client

	id        int
	firstName string
	username  string

	offset   *int
	offsetMu *sync.RWMutex

	admins   set.Set[chatdomain.UserId]
	batch    mexadomain.Batch
	exercise mexadomain.Exercise
}

func (bot *Bot) Exercise() mexadomain.Exercise {
	return bot.exercise
}

func (bot *Bot) Batch() mexadomain.Batch {
	return bot.batch
}

func New(ctx context.Context, token string, batch mexadomain.Batch, exercise mexadomain.Exercise, admins []chatdomain.UserId) (bot *Bot, err error) {
	bot = &Bot{
		token:   token,
		baseUrl: telegramBaseUrl + token,
		client:  http.DefaultClient,

		offset:   nil,
		offsetMu: &sync.RWMutex{},

		batch:    batch,
		exercise: exercise,
	}

	bot.admins = set.New(set.WithSlice(admins))

	err = bot.getMe(ctx)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (bot *Bot) getMe(ctx context.Context) (err error) {
	url := bot.baseUrl + "/getMe"
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := bot.client.Do(r)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var body struct {
		Ok     bool `json:"ok"`
		Result struct {
			Id        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"result"`
	}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}

	bot.id = body.Result.Id
	bot.firstName = body.Result.FirstName
	bot.username = body.Result.Username
	return nil
}

func (bot *Bot) UpdateOffset(offset int) {
	offset++

	bot.offsetMu.Lock()
	defer bot.offsetMu.Unlock()

	if bot.offset == nil || offset >= *bot.offset {
		bot.offset = &offset
	}
}

func (bot *Bot) GetUpdates() (u []chatdomain.Update, err error) {
	url := bot.baseUrl + "/getUpdates"
	if bot.offset != nil {
		url += fmt.Sprintf("?offset=%d", *bot.offset)
	}

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := bot.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("BODY", string(b))

	var res struct {
		Ok     bool                `json:"ok"`
		Result []chatdomain.Update `json:"result"`
	}

	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	} else if !res.Ok {
		fmt.Println(string(b))
		return nil, fmt.Errorf("getUpdates: not ok")
	}

	return res.Result, nil
}

func (bot *Bot) Reply(ctx context.Context, chatId chatdomain.ChatId, text string, opts ...chatdomain.SendOpt) (err error) {
	url := bot.baseUrl + "/sendMessage"

	m := map[string]any{
		"chat_id": chatId,
		"text":    text,
	}

	conf := defaultSendConfig()
	for _, opt := range opts {
		opt(&conf)
	}

	m["parse_mode"] = conf.ParseMode
	if conf.Rmk != nil {
		m["reply_markup"] = *conf.Rmk
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	fmt.Println("BODY", string(b))
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")

	var body struct {
		Ok bool `json:"ok"`
	}
	resp, err := bot.client.Do(r)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading error body", err)
		} else {
			fmt.Println("ERROR BODY", string(b))
		}

		return fmt.Errorf("reply: unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	} else if !body.Ok {
		return fmt.Errorf("reply: not ok")
	}

	return nil
}

func defaultSendConfig() chatdomain.SendConfig {
	return chatdomain.SendConfig{
		Rmk:       nil,
		ParseMode: "MarkdownV2",
	}
}

func (bot *Bot) EditMessage(ctx context.Context, chatId chatdomain.ChatId, messageId chatdomain.MessageId, text string, opts ...chatdomain.SendOpt) (err error) {
	url := bot.baseUrl + "/editMessageText"
	m := map[string]any{
		"chat_id":    chatId,
		"message_id": messageId,
		"text":       text,
	}

	conf := defaultSendConfig()
	for _, opt := range opts {
		opt(&conf)
	}

	m["parse_mode"] = conf.ParseMode
	if conf.Rmk != nil {
		m["reply_markup"] = *conf.Rmk
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := bot.client.Do(r)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("editMessage: unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var body struct {
		Ok bool `json:"ok"`
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		return err
	} else if !body.Ok {
		return fmt.Errorf("editMessage: not ok")
	}

	return nil
}
