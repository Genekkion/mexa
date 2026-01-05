package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	"net/http"
)

func (bot *Bot) SetupCommands(ctx context.Context, commands []chatdomain.Command) (err error) {
	if len(commands) == 0 {
		return nil
	}

	body := map[string]any{
		"commands": commands,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := bot.baseUrl + "/setMyCommands"
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")

	var respBody struct {
		Ok bool `json:"ok"`
	}

	resp, err := bot.client.Do(r)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("setMyCommands: unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	} else if !respBody.Ok {
		return fmt.Errorf("setMyCommands: not ok")
	}

	return nil
}
