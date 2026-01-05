package worker

import (
	"context"
	"fmt"
	chatdomain "mexa/internal/domains/chat"
	botports "mexa/internal/ports/bot"
	mexaservice "mexa/internal/services/mexa"
	"sync"
	"time"
)

type Worker struct {
	bot botports.Bot
	ser mexaservice.Service
}

func New(bot botports.Bot, ser mexaservice.Service) *Worker {
	return &Worker{
		bot: bot,
		ser: ser,
	}
}

func (w *Worker) Start(ctx context.Context) (err error) {
	timer := time.NewTimer(time.Second)

	wg := sync.WaitGroup{}
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return nil
		case <-timer.C:
			var updates []chatdomain.Update
			updates, err = w.bot.GetUpdates()
			if err != nil {
				return err
			}

			for _, u := range updates {
				wg.Go(func() {
					w.bot.UpdateOffset(u.UpdateId)

					fmt.Printf("Handling update: %+v\n", u)
					err := w.handleUpdate(ctx, u)
					if err != nil {
						fmt.Printf("Error handling update: %v\n", err)
					}
				})
			}

			timer.Reset(time.Second)
		}
	}
}

func (w *Worker) handleUpdate(ctx context.Context, u chatdomain.Update) (err error) {
	if u.ToIgnore() {
		return nil
	} else if u.IsCommand() {
		return w.ser.HandleCommands(ctx, u)
	} else if u.IsCallback() {
		return w.ser.HandleCallbacks(ctx, u)
	}

	return w.ser.HandleText(ctx, u)
}
