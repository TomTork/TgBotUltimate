package telegram

import (
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"os"
)

func Telegram(ctx context.Context, botToken string) error {
	bot, err := telego.NewBot(botToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update, ok := <-updates:
			if !ok || update.Message == nil {
				return nil
			}
			if update.Message != nil {
				_, err := bot.SendMessage(ctx, tu.Message(tu.ID(update.Message.Chat.ID), update.Message.Text))
				if err != nil {
					return err
				}
			}
		}
	}
}
