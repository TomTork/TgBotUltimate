package actions

import (
	"TgBotUltimate/database/data"
	"TgBotUltimate/database/favorites"
	"TgBotUltimate/database/messages"
	"TgBotUltimate/database/users"
	"TgBotUltimate/processing"
	"TgBotUltimate/processing/neuro"
	"TgBotUltimate/types/Action"
	"TgBotUltimate/types/Database"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"strings"
)

const (
	ShowMoreFlatsPrefix      = "show_more_flats"
	ShowFavoriteFlatsPrefix  = "show_favorite_flats"
	ToggleFavoriteFlatPrefix = "toggle_favorite_flat"
)

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

func ToggleFavoriteFlat(a Action.Action) error {
	if a.Update.CallbackQuery == nil || a.Update.CallbackQuery.Message == nil {
		return nil
	}

	callback := a.Update.CallbackQuery
	if err := answerShowMoreCallback(a); err != nil {
		return err
	}
	if err := ensureCallbackUser(a); err != nil {
		return err
	}

	parts := strings.Split(callback.Data, ":")
	if len(parts) != 3 {
		return nil
	}

	flatCode := parts[1]
	desiredState := parts[2]

	switch desiredState {
	case "add":
		if err := favorites.AddFavorite(a.ReqCtx, a.Database, callback.From.ID, flatCode); err != nil {
			return err
		}
	case "remove":
		if err := favorites.RemoveFavorite(a.ReqCtx, a.Database, callback.From.ID, flatCode); err != nil {
			return err
		}
	default:
		return nil
	}

	_, err := a.Bot.EditMessageReplyMarkup(a.ReqCtx, &telego.EditMessageReplyMarkupParams{
		ChatID:      tu.ID(callback.Message.GetChat().ID),
		MessageID:   callback.Message.GetMessageID(),
		ReplyMarkup: favoriteFlatKeyboard(flatCode, desiredState == "add"),
	})
	return err
}

func ShowFavoriteFlats(a Action.Action) error {
	var userID int64
	var chatID int64

	if a.Update.CallbackQuery != nil && a.Update.CallbackQuery.Message != nil {
		if err := answerShowMoreCallback(a); err != nil {
			return err
		}
		if err := ensureCallbackUser(a); err != nil {
			return err
		}
		userID = a.Update.CallbackQuery.From.ID
		chatID = a.Update.CallbackQuery.Message.GetChat().ID
	} else if a.Update.Message != nil && a.Update.Message.From != nil {
		userID = a.Update.Message.From.ID
		chatID = a.Update.Message.Chat.ID
	} else {
		return nil
	}

	flats, err := favorites.GetFavoriteFlatsByUser(a.ReqCtx, a.Database, userID)
	if err != nil {
		return err
	}
	if len(flats) == 0 {
		_, err = a.Bot.SendMessage(
			a.Ctx,
			tu.Message(tu.ID(chatID), "В избранном пока нет планировок."),
		)
		return err
	}

	return sendFlatCards(a, userID, chatID, flats, true, nil)
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

	if err := sendFlatCards(a, *user.TgId, chatID, flats, false, user); err != nil {
		return err
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

func sendFlatCards(a Action.Action, userID int64, chatID int64, flats []Database.Query, favoritesOnly bool, user *Database.User) error {
	sentFlats := 0
	for _, flat := range flats {
		show, flatImg, _ := processing.ShowFlat(flat)
		if flatImg == "" {
			continue
		}
		if flat.FlatCode == nil || *flat.FlatCode == "" {
			continue
		}

		isFavorite, err := favorites.IsFavorite(a.ReqCtx, a.Database, userID, *flat.FlatCode)
		if err != nil {
			return err
		}

		_, err = a.Bot.SendPhoto(
			a.Ctx,
			tu.Photo(
				tu.ID(chatID),
				tu.FileFromURL(flatImg),
			).WithCaption(show).WithReplyMarkup(favoriteFlatKeyboard(*flat.FlatCode, isFavorite)),
		)
		if err != nil {
			return err
		}
		sentFlats++
	}

	if sentFlats == 0 {
		if favoritesOnly {
			_, err := a.Bot.SendMessage(
				a.Ctx,
				tu.Message(tu.ID(chatID), "В избранном нет доступных планировок с изображением."),
			)
			return err
		}
		return sendNoFlatsFound(a, user, chatID)
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
			tu.InlineKeyboardButton("Показать избранное").WithCallbackData(ShowFavoriteFlatsPrefix),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Сбросить параметры").WithCallbackData(ExpertResetPrefix),
		),
	)
}

func favoriteFlatKeyboard(flatCode string, isFavorite bool) *telego.InlineKeyboardMarkup {
	if isFavorite {
		return tu.InlineKeyboard(
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("Удалить из избранного").WithCallbackData(ToggleFavoriteFlatPrefix + ":" + flatCode + ":remove"),
			),
		)
	}

	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Добавить в избранное").WithCallbackData(ToggleFavoriteFlatPrefix + ":" + flatCode + ":add"),
		),
	)
}

func answerShowMoreCallback(a Action.Action) error {
	return a.Bot.AnswerCallbackQuery(a.ReqCtx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: a.Update.CallbackQuery.ID,
	})
}
