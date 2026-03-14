package telegram

import (
	"TgBotUltimate/database/data"
	"TgBotUltimate/database/messages"
	"TgBotUltimate/database/users"
	"TgBotUltimate/processing"
	"TgBotUltimate/processing/neuro"
	"TgBotUltimate/types/Database"
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Telegram(ctx context.Context, botToken string, database *Database.DB) error {
	bot, err := telego.NewBot(botToken)
	if err != nil {
		return fmt.Errorf("create telegram bot: %w", err)
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
			if update.Message == nil || update.Message.From == nil {
				continue
			}
			reqCtx := context.WithoutCancel(ctx)
			_ = users.CreateUser(
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
			n, err := neuro.Parameters(ctx, update.Message.Text)
			if err != nil {
				return fmt.Errorf("parse neuro parameters: %w", err)
			}
			_ = messages.CreateMessage(
				reqCtx,
				database,
				Database.ChatMessage{
					TgId:    uint64(update.Message.From.ID),
					Message: update.Message.Text,
					Parameters: Database.Parameters{
						ProjectName:    string(n.ProjectName),
						BuildingLiter:  string(n.BuildingLiter),
						FloorMin:       string(n.FloorMin),
						FloorMax:       string(n.FloorMax),
						RoomsAmountMin: string(n.RoomsAmountMin),
						RoomsAmountMax: string(n.RoomsAmountMax),
						SquareMin:      string(n.SquareMin),
						SquareMax:      string(n.SquareMax),
						CostMin:        string(n.CostMin),
						CostMax:        string(n.CostMax),
					},
				})
			user, err := users.GetUserById(reqCtx, database, update.Message.From.ID)
			if err != nil {
				return fmt.Errorf("get user: %w", err)
			}
			flats, err := data.GetFlatsByParameters(reqCtx, database, user)
			for _, flat := range flats {
				show, flatImg, _ := processing.ShowFlat(flat)
				if flatImg != "" {
					_, err = bot.SendPhoto(
						ctx,
						tu.Photo(
							tu.ID(update.Message.Chat.ID),
							tu.FileFromURL(flatImg),
						).WithCaption(show),
					)
				}
				//_, err = bot.SendMessage(
				//	ctx,
				//	tu.Message(
				//		tu.ID(update.Message.Chat.ID),
				//		show,
				//	),
				//)
			}
			if err != nil {
				return err
			}
		}
	}
}
