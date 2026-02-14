package telegram

import (
	"TgBotUltimate/database/messages"
	"TgBotUltimate/database/users"
	"TgBotUltimate/processing/neuro"
	"TgBotUltimate/types/Database"
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"os"
	"strings"
)

func Telegram(ctx context.Context, botToken string, database *Database.DB) error {
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
			if !ok {
				return nil
			}
			if update.Message == nil {
				continue
			}
			reqCtx := context.WithoutCancel(ctx)
			err = users.CreateUser(
				reqCtx,
				database,
				Database.User{
					TgId:        &update.Message.From.ID,
					UserName:    &update.Message.From.Username,
					FirstName:   &update.Message.From.FirstName,
					LastName:    &update.Message.From.LastName,
					PhoneNumber: nil,
					Email:       nil,
				})
			err = messages.CreateMessage(reqCtx, database, Database.ChatMessage{TgId: uint64(update.Message.From.ID), Message: update.Message.Text})
			if err != nil {
				return err
			}
			_messages, _ := messages.GetMessagesByTgId(reqCtx, database, uint64(update.Message.From.ID))
			__messages := make([]string, 0, len(_messages))
			for _, message := range _messages {
				__messages = append(__messages, message.Message)
			}
			_, err = bot.SendMessage(
				ctx,
				tu.Message(
					tu.ID(update.Message.Chat.ID),
					fmt.Sprintf("%d"+
						"\nТвой текст запроса:\n%s"+
						"\nПредыдущие запросы:\n\n%s"+
						"\nВыделенные переменные:\n\n%s",
						update.Message.From.ID,
						update.Message.Text,
						strings.Join(__messages, "\n"),
						neuro.Parameters(ctx, update.Message.Text),
					),
				))
			if err != nil {
				return err
			}
		}
	}
}
