package telegram

import (
	"TgBotUltimate/database/users"
	"TgBotUltimate/platform/actions"
	"TgBotUltimate/types/Action"
	"TgBotUltimate/types/Database"
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	"log"
)

func Telegram(ctx context.Context, botToken string, database *Database.DB) error {
	bot, err := telego.NewBot(botToken)
	if err != nil {
		return fmt.Errorf("create telegram bot: %w", err)
	}

	// Установка команд
	actions.SetCommands(Action.Action{
		Ctx: context.WithoutCancel(ctx),
		Bot: bot,
	})

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
			reqCtx := context.WithoutCancel(ctx)
			action := Action.Action{
				ReqCtx:   reqCtx,
				Ctx:      ctx,
				Update:   update,
				Database: database,
				Bot:      bot,
			}

			if update.CallbackQuery != nil {
				actions.CallbackQuery(action)
				continue
			}

			if update.Message == nil || update.Message.From == nil {
				continue
			}

			_ = users.CreateUser(
				action.ReqCtx,
				action.Database,
				Database.User{
					TgId:        &action.Update.Message.From.ID,
					UserName:    &action.Update.Message.From.Username,
					FirstName:   &action.Update.Message.From.FirstName,
					LastName:    &action.Update.Message.From.LastName,
					PhoneNumber: nil,
					Email:       nil,
				})

			switch update.Message.Text {
			case "/start":
				start := actions.Start(action)
				if start != nil {
					log.Println(start)
					return start
				}
			case "/help":
				help := actions.Help(action)
				if help != nil {
					log.Println(help)
					return help
				}
			case "/questions":
				questions := actions.StartExpertSystemCommand(action)
				if questions != nil {
					log.Println(questions)
					return questions
				}
			case "/reload":
				reload := actions.ResetExpertSystemCommand(action)
				if reload != nil {
					log.Println(reload)
					return reload
				}
			case "/flats":
				flats := actions.StartFlatSelectionCommand(action)
				if flats != nil {
					log.Println(flats)
					return flats
				}
			default:
				handled, err := actions.HandleManualParameterMessage(action)
				if err != nil {
					log.Println(err)
					return err
				}
				if handled {
					continue
				}
				selection := actions.Selection(action)
				if selection != nil {
					log.Println(selection)
					return selection
				}
			}
		}
	}
}
