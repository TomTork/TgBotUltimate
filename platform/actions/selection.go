package actions

import (
	"TgBotUltimate/database/data"
	"TgBotUltimate/database/messages"
	"TgBotUltimate/database/users"
	"TgBotUltimate/processing"
	"TgBotUltimate/processing/neuro"
	"TgBotUltimate/types/Action"
	"TgBotUltimate/types/Database"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const ShowMoreFlatsPrefix = "show_more_flats"

func Selection(a Action.Action) error {
	n, err := neuro.Parameters(a.Ctx, a.Update.Message.Text)
	if err != nil {
		return fmt.Errorf("parse neuro parameters: %w", err)
	}
	_ = messages.CreateMessage(
		a.ReqCtx,
		a.Database,
		Database.ChatMessage{
			TgId:    uint64(a.Update.Message.From.ID),
			Message: a.Update.Message.Text,
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
	user, err := users.GetUserById(a.ReqCtx, a.Database, a.Update.Message.From.ID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	return sendFlatsByUser(a, user, a.Update.Message.Chat.ID, false)
}

func ShowMoreFlats(a Action.Action) error {
	if a.Update.CallbackQuery == nil || a.Update.CallbackQuery.Message == nil {
		return nil
	}

	if err := answerShowMoreCallback(a); err != nil {
		return err
	}

	user, err := users.GetUserById(a.ReqCtx, a.Database, a.Update.CallbackQuery.From.ID)
	if err != nil {
		return err
	}

	return sendFlatsByUser(a, user, a.Update.CallbackQuery.Message.GetChat().ID, true)
}

func sendFlatsByUser(a Action.Action, user *Database.User, chatID int64, increaseOffset bool) error {
	if user == nil {
		return nil
	}

	flats, err := data.GetFlatsByParameters(a.ReqCtx, a.Database, user)
	if err != nil {
		return err
	}
	if len(flats) == 0 {
		return sendNoFlatsFound(a, user, chatID)
	}

	sentFlats := 0
	for _, flat := range flats {
		show, flatImg, _ := processing.ShowFlat(flat)
		if flatImg == "" {
			continue
		}

		_, err = a.Bot.SendPhoto(
			a.Ctx,
			tu.Photo(
				tu.ID(chatID),
				tu.FileFromURL(flatImg),
			).WithCaption(show),
		)
		if err != nil {
			return err
		}
		sentFlats++
	}

	if sentFlats == 0 {
		return sendNoFlatsFound(a, user, chatID)
	}

	_, err = a.Bot.SendMessage(
		a.Ctx,
		tu.Message(
			tu.ID(chatID),
			"Смотрим следующие квартиры?",
		).WithReplyMarkup(showMoreKeyboard()),
	)
	if err != nil {
		return err
	}

	if increaseOffset {
		return users.IncreaseUserOffset(a.Ctx, a.Database, *user.TgId)
	}

	return nil
}

func sendNoFlatsFound(a Action.Action, user *Database.User, chatID int64) error {
	if user != nil && user.TgId != nil {
		_ = users.DropUserOffset(a.Ctx, a.Database, *user.TgId)
	}
	_, err := a.Bot.SendMessage(
		a.Ctx,
		tu.Message(
			tu.ID(chatID),
			"Квартиры не найдены.",
		).WithReplyMarkup(tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("Сбросить варианты ответов").WithCallbackData(ExpertResetPrefix),
			),
		)),
	)
	return err
}

func showMoreKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Показать ещё").WithCallbackData(ShowMoreFlatsPrefix),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Сбросить параметры").WithCallbackData(ExpertResetPrefix),
		),
	)
}

func answerShowMoreCallback(a Action.Action) error {
	return a.Bot.AnswerCallbackQuery(a.ReqCtx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: a.Update.CallbackQuery.ID,
	})
}
