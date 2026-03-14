package actions

import (
	"TgBotUltimate/database/expert"
	"TgBotUltimate/database/users"
	"TgBotUltimate/types/Action"
	dbtypes "TgBotUltimate/types/Database"
	"TgBotUltimate/types/Expert"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
	"strings"
)

const (
	ExpertStartPrefix       = "expert_system"
	ExpertAnswerPrefix      = "expert_answer"
	ExpertFinishPrefix      = "expert_finish"
	ExpertResetPrefix       = "expert_reset"
	ExpertSelectFlatsPrefix = "expert_select_flats"
)

func ExpertSystem(action Action.Action) error {
	questions, err := expert.GetQuestions(action.Ctx, action.Database)
	if err != nil {
		return err
	}
	if len(questions) == 0 {
		return sendText(action, "В экспертной системе пока нет вопросов.")
	}
	if action.Update.CallbackQuery == nil {
		return nil
	}
	callback := action.Update.CallbackQuery
	if err := answerCallback(action, callback.ID); err != nil {
		return err
	}
	if callback.Message == nil {
		return nil
	}
	data := callback.Data
	switch {
	case data == ExpertStartPrefix:
		return showQuestion(action, 0, false)
	case data == ExpertFinishPrefix:
		return finishExpertSystem(action)
	case data == ExpertResetPrefix:
		return resetExpertSystem(action)
	case data == ExpertSelectFlatsPrefix:
		return startFlatSelection(action)
	case strings.HasPrefix(data, ExpertAnswerPrefix+":"):
		return handleAnswerCallback(action, questions, data)
	default:
		return nil
	}
}

func handleAnswerCallback(action Action.Action, questions []Expert.Question, data string) error {
	// ожидаемый формат:
	// expert_answer:<questionIndex>:<variantIndex>
	parts := strings.Split(data, ":")
	if len(parts) != 3 {
		return nil
	}

	questionIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	variantIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	if questionIndex < 0 || questionIndex >= len(questions) {
		return finishExpertSystem(action)
	}

	currentQuestion := questions[questionIndex]
	variants := splitAndTrim(currentQuestion.Variants)

	if variantIndex < 0 || variantIndex >= len(variants) {
		return nil
	}

	err = processExpertAnswer(action, currentQuestion.Results, variantIndex)
	if err != nil {
		return err
	}

	nextIndex := questionIndex + 1
	if nextIndex >= len(questions) {
		return finishExpertSystem(action)
	}

	return showQuestion(action, nextIndex, true)
}

func showQuestion(action Action.Action, questionIndex int, editCurrentMessage bool) error {
	questions, err := expert.GetQuestions(action.Ctx, action.Database)
	if err != nil {
		return err
	}
	if questionIndex < 0 || questionIndex >= len(questions) {
		return finishExpertSystem(action)
	}

	question := questions[questionIndex]
	variants := splitAndTrim(question.Variants)

	buttonRows := make([][]telego.InlineKeyboardButton, 0, len(variants)+1)

	for i, variant := range variants {
		buttonRows = append(buttonRows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(variant).
				WithCallbackData(fmt.Sprintf("%s:%d:%d", ExpertAnswerPrefix, questionIndex, i)),
		))
	}

	// обязательная кнопка завершения
	buttonRows = append(buttonRows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("Завершить").WithCallbackData(ExpertFinishPrefix),
	))

	keyboard := tu.InlineKeyboard(buttonRows...)

	text := fmt.Sprintf(
		"Вопрос %d из %d\n\n%s",
		questionIndex+1,
		len(questions),
		question.Question,
	)

	callback := action.Update.CallbackQuery
	if callback == nil || callback.Message == nil {
		return nil
	}

	chatID := callback.Message.GetChat().ID
	messageID := callback.Message.GetMessageID()

	if editCurrentMessage {
		_, err := action.Bot.EditMessageText(action.ReqCtx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(chatID),
			MessageID:   messageID,
			Text:        text,
			ReplyMarkup: keyboard,
		})
		return err
	}

	_, err = action.Bot.SendMessage(action.ReqCtx, tu.Message(
		tu.ID(chatID),
		text,
	).WithReplyMarkup(keyboard))
	return err
}

func finishExpertSystem(action Action.Action) error {
	callback := action.Update.CallbackQuery
	if callback == nil || callback.Message == nil {
		return nil
	}

	_, _ = action.Bot.EditMessageText(action.ReqCtx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(callback.Message.GetChat().ID),
		MessageID: callback.Message.GetMessageID(),
		Text:      "Серия вопросов завершена.",
	})

	_, err := action.Bot.SendMessage(action.ReqCtx, tu.Message(
		tu.ID(callback.Message.GetChat().ID),
		"Серия вопросов завершена.",
	).WithReplyMarkup(expertFinishKeyboard()))
	return err
}

func answerCallback(action Action.Action, callbackID string) error {
	return action.Bot.AnswerCallbackQuery(action.ReqCtx, &telego.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
	})
}

func sendText(action Action.Action, text string) error {
	var chatID int64

	if action.Update.CallbackQuery != nil && action.Update.CallbackQuery.Message != nil {
		chatID = action.Update.CallbackQuery.Message.GetChat().ID
	} else if action.Update.Message != nil {
		chatID = action.Update.Message.Chat.ID
	} else {
		return nil
	}

	_, err := action.Bot.SendMessage(context.Background(), tu.Message(
		tu.ID(chatID),
		text,
	))
	return err
}

func splitAndTrim(value string) []string {
	raw := strings.Split(value, ",")
	result := make([]string, 0, len(raw))

	for _, item := range raw {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}

func processExpertAnswer(action Action.Action, results string, variantIndex int) error {
	if strings.TrimSpace(results) == "" {
		return nil
	}

	callback := action.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	if err := ensureCallbackUser(action); err != nil {
		return err
	}

	var parsedResults map[string]map[string]string
	if err := json.Unmarshal([]byte(results), &parsedResults); err != nil {
		return nil
	}

	variantResult, ok := parsedResults[strconv.Itoa(variantIndex+1)]
	if !ok {
		return nil
	}

	system := buildExpertSystemUpdate(variantResult)
	return users.SetExpertSystemFields(action.Ctx, action.Database, callback.From.ID, system)
}

func resetExpertSystem(action Action.Action) error {
	callback := action.Update.CallbackQuery
	if callback == nil || callback.Message == nil {
		return nil
	}

	if err := ensureCallbackUser(action); err != nil {
		return err
	}
	if err := users.ResetExpertSystemFields(action.Ctx, action.Database, callback.From.ID); err != nil {
		return err
	}

	_, err := action.Bot.EditMessageText(action.ReqCtx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(callback.Message.GetChat().ID),
		MessageID:   callback.Message.GetMessageID(),
		Text:        "Параметры expert system сброшены.",
		ReplyMarkup: expertFinishKeyboard(),
	})
	return err
}

func startFlatSelection(action Action.Action) error {
	callback := action.Update.CallbackQuery
	if callback == nil || callback.Message == nil {
		return nil
	}

	if err := ensureCallbackUser(action); err != nil {
		return err
	}

	user, err := users.GetUserById(action.Ctx, action.Database, callback.From.ID)
	if err != nil {
		return err
	}

	return sendFlatsByUser(action, user, callback.Message.GetChat().ID, true)
}

func buildExpertSystemUpdate(values map[string]string) dbtypes.ExpertSystem {
	system := dbtypes.ExpertSystem{}

	for key, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		switch key {
		case "project_name":
			system.ExProjectName = stringPtr(value)
		case "building_liter":
			system.ExBuildingLiter = stringPtr(value)
		case "floor_min":
			system.ExFloorMin = stringPtr(value)
		case "floor_max":
			system.ExFloorMax = stringPtr(value)
		case "rooms_amount_min":
			system.ExRoomsAmountMin = stringPtr(value)
		case "rooms_amount_max":
			system.ExRoomsAmountMax = stringPtr(value)
		case "square_min":
			system.ExSquareMin = stringPtr(value)
		case "square_max":
			system.ExSquareMax = stringPtr(value)
		case "cost_min":
			system.ExCostMin = stringPtr(value)
		case "cost_max":
			system.ExCostMax = stringPtr(value)
		}
	}

	return system
}

func stringPtr(value string) *string {
	return &value
}

func expertFinishKeyboard() *telego.InlineKeyboardMarkup {
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Сбросить варианты ответов").WithCallbackData(ExpertResetPrefix),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Начать подбор квартиры").WithCallbackData(ExpertSelectFlatsPrefix),
		),
	)
}

func ensureCallbackUser(action Action.Action) error {
	callback := action.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	return users.CreateUser(
		action.ReqCtx,
		action.Database,
		dbtypes.User{
			TgId:        &callback.From.ID,
			UserName:    &callback.From.Username,
			FirstName:   &callback.From.FirstName,
			LastName:    &callback.From.LastName,
			PhoneNumber: nil,
			Email:       nil,
		},
	)
}
