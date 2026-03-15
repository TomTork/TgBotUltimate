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
	ExpertNextPrefix        = "expert_next"
	ExpertFinishPrefix      = "expert_finish"
	ExpertResetPrefix       = "expert_reset"
	ExpertSelectFlatsPrefix = "expert_select_flats"
)

func ExpertSystem(action Action.Action) error {
	questions, err := expert.GetQuestions(action.Ctx, action.Database)
	if err != nil {
		return err
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
		deleteManualParameterState(callback.From.ID)
		return showNextQuestion(action, questions, nil, false)
	case data == ExpertFinishPrefix:
		return finishExpertSystem(action)
	case data == ExpertResetPrefix:
		return resetExpertSystem(action)
	case data == ExpertSelectFlatsPrefix:
		return startFlatSelection(action)
	case strings.HasPrefix(data, ExpertAnswerPrefix+":"):
		return handleAnswerCallback(action, questions, data)
	case strings.HasPrefix(data, ExpertNextPrefix+":"):
		return handleSkipCallback(action, questions, data)
	default:
		return nil
	}
}

func handleAnswerCallback(action Action.Action, questions []Expert.Question, data string) error {
	// expert_answer:<questionID>:<variantIndex>
	parts := strings.Split(data, ":")
	if len(parts) != 3 {
		return nil
	}

	questionID, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	variantIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	currentQuestion, ok := findQuestionByID(questions, questionID)
	if !ok {
		return finishExpertSystem(action)
	}
	available, err := isQuestionAvailable(action, questions, questionID)
	if err != nil {
		return err
	}
	if !available {
		return showNextQuestion(action, questions, &questionID, true)
	}
	variants := splitAndTrim(currentQuestion.Variants)

	if variantIndex < 0 || variantIndex >= len(variants) {
		return nil
	}

	err = processExpertAnswer(action, currentQuestion, variantIndex)
	if err != nil {
		return err
	}

	return showNextQuestion(action, questions, &questionID, true)
}

func handleSkipCallback(action Action.Action, questions []Expert.Question, data string) error {
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return nil
	}

	questionID, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}

	return showNextQuestion(action, questions, &questionID, true)
}

func showNextQuestion(action Action.Action, questions []Expert.Question, afterQuestionID *int, editCurrentMessage bool) error {
	nextQuestion, displayIndex, totalAvailable, err := getNextQuestion(action, questions, afterQuestionID)
	if err != nil {
		return err
	}
	if nextQuestion == nil {
		return finishExpertSystem(action)
	}

	text, keyboard := buildExpertQuestionView(action, questions, *nextQuestion, displayIndex, totalAvailable)

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

func StartExpertSystemCommand(action Action.Action) error {
	if action.Update.Message == nil || action.Update.Message.From == nil {
		return nil
	}

	questions, err := expert.GetQuestions(action.Ctx, action.Database)
	if err != nil {
		return err
	}

	deleteManualParameterState(action.Update.Message.From.ID)

	nextQuestion, displayIndex, totalAvailable, err := getNextQuestion(action, questions, nil)
	if err != nil {
		return err
	}
	if nextQuestion == nil {
		_, err = action.Bot.SendMessage(action.ReqCtx, tu.Message(
			tu.ID(action.Update.Message.Chat.ID),
			"В экспертной системе пока нет доступных вопросов.",
		))
		return err
	}

	text, keyboard := buildExpertQuestionView(action, questions, *nextQuestion, displayIndex, totalAvailable)
	_, err = action.Bot.SendMessage(action.ReqCtx, tu.Message(
		tu.ID(action.Update.Message.Chat.ID),
		text,
	).WithReplyMarkup(keyboard))
	return err
}

func buildExpertQuestionView(action Action.Action, questions []Expert.Question, question Expert.Question, displayIndex int, totalAvailable int) (string, *telego.InlineKeyboardMarkup) {
	variants := splitAndTrim(question.Variants)

	buttonRows := make([][]telego.InlineKeyboardButton, 0, len(variants)+1)

	for i, variant := range variants {
		buttonRows = append(buttonRows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(variant).
				WithCallbackData(fmt.Sprintf("%s:%d:%d", ExpertAnswerPrefix, question.Id, i)),
		))
	}

	navigationRow := make([]telego.InlineKeyboardButton, 0, 2)
	if hasNextQuestion(action, questions, question.Id) {
		navigationRow = append(navigationRow,
			tu.InlineKeyboardButton("▶").WithCallbackData(fmt.Sprintf("%s:%d", ExpertNextPrefix, question.Id)),
		)
	}
	navigationRow = append(navigationRow,
		tu.InlineKeyboardButton("Завершить").WithCallbackData(ExpertFinishPrefix),
	)

	buttonRows = append(buttonRows, navigationRow)

	text := fmt.Sprintf(
		"Вопрос %d из %d\n\n%s",
		displayIndex,
		totalAvailable,
		question.Question,
	)

	return text, tu.InlineKeyboard(buttonRows...)
}

func getNextQuestion(action Action.Action, questions []Expert.Question, afterQuestionID *int) (*Expert.Question, int, int, error) {
	answers, err := getCurrentUserExpertAnswers(action)
	if err != nil {
		return nil, 0, 0, err
	}

	excludedQuestionIDs := buildExcludedQuestionIDs(questions, answers)
	availableQuestions := filterAvailableQuestions(questions, excludedQuestionIDs)
	if len(availableQuestions) == 0 {
		return nil, 0, 0, nil
	}

	startIndex := 0
	if afterQuestionID != nil {
		startIndex = len(questions)
		for i, question := range questions {
			if question.Id == *afterQuestionID {
				startIndex = i + 1
				break
			}
		}
	}

	for i := startIndex; i < len(questions); i++ {
		question := questions[i]
		if _, excluded := excludedQuestionIDs[question.Id]; excluded {
			continue
		}

		return &question, questionPosition(availableQuestions, question.Id), len(availableQuestions), nil
	}

	return nil, 0, len(availableQuestions), nil
}

func hasNextQuestion(action Action.Action, questions []Expert.Question, currentQuestionID int) bool {
	nextQuestion, _, _, err := getNextQuestion(action, questions, &currentQuestionID)
	return err == nil && nextQuestion != nil
}

func isQuestionAvailable(action Action.Action, questions []Expert.Question, questionID int) (bool, error) {
	answers, err := getCurrentUserExpertAnswers(action)
	if err != nil {
		return false, err
	}

	excludedQuestionIDs := buildExcludedQuestionIDs(questions, answers)
	_, excluded := excludedQuestionIDs[questionID]
	return !excluded, nil
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

func processExpertAnswer(action Action.Action, question Expert.Question, variantIndex int) error {
	callback := action.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	if err := ensureCallbackUser(action); err != nil {
		return err
	}

	if err := users.SaveExpertSystemAnswer(action.Ctx, action.Database, dbtypes.ExpertSystemAnswer{
		UserTgID:     callback.From.ID,
		QuestionID:   question.Id,
		VariantIndex: variantIndex,
	}); err != nil {
		return err
	}

	if strings.TrimSpace(question.Results) == "" {
		return nil
	}

	var parsedResults map[string]map[string]string
	if err := json.Unmarshal([]byte(question.Results), &parsedResults); err != nil {
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
	if err := resetExpertSystemByUserID(action, callback.From.ID); err != nil {
		return err
	}

	_, err := action.Bot.EditMessageText(action.ReqCtx, &telego.EditMessageTextParams{
		ChatID:      tu.ID(callback.Message.GetChat().ID),
		MessageID:   callback.Message.GetMessageID(),
		Text:        "Параметры экспертной системы сброшены.",
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
	return startFlatSelectionByUserID(action, callback.From.ID, callback.Message.GetChat().ID, true)
}

func ResetExpertSystemCommand(action Action.Action) error {
	if action.Update.Message == nil || action.Update.Message.From == nil {
		return nil
	}

	if err := resetExpertSystemByUserID(action, action.Update.Message.From.ID); err != nil {
		return err
	}

	_, err := action.Bot.SendMessage(action.ReqCtx, tu.Message(
		tu.ID(action.Update.Message.Chat.ID),
		"Параметры экспертной системы сброшены.",
	).WithReplyMarkup(expertFinishKeyboard()))
	return err
}

func StartFlatSelectionCommand(action Action.Action) error {
	if action.Update.Message == nil || action.Update.Message.From == nil {
		return nil
	}

	return startFlatSelectionByUserID(action, action.Update.Message.From.ID, action.Update.Message.Chat.ID, false)
}

func resetExpertSystemByUserID(action Action.Action, userID int64) error {
	if err := users.ResetExpertSystemAnswers(action.Ctx, action.Database, userID); err != nil {
		return err
	}
	if err := users.ResetExpertSystemFields(action.Ctx, action.Database, userID); err != nil {
		return err
	}
	deleteManualParameterState(userID)
	return nil
}

func startFlatSelectionByUserID(action Action.Action, userID int64, chatID int64, increaseOffset bool) error {
	deleteManualParameterState(userID)

	user, err := users.GetUserById(action.Ctx, action.Database, userID)
	if err != nil {
		return err
	}

	return sendFlatsByUser(action, user, chatID, increaseOffset)
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
			tu.InlineKeyboardButton("Сбросить вопросы экспертной системы").WithCallbackData(ExpertResetPrefix),
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

func getCurrentUserExpertAnswers(action Action.Action) ([]dbtypes.ExpertSystemAnswer, error) {
	if action.Update.CallbackQuery != nil {
		callback := action.Update.CallbackQuery
		if err := ensureCallbackUser(action); err != nil {
			return nil, err
		}
		return users.GetExpertSystemAnswers(action.Ctx, action.Database, callback.From.ID)
	}

	if action.Update.Message != nil && action.Update.Message.From != nil {
		return users.GetExpertSystemAnswers(action.Ctx, action.Database, action.Update.Message.From.ID)
	}

	return nil, nil
}

func filterAvailableQuestions(questions []Expert.Question, excludedQuestionIDs map[int]struct{}) []Expert.Question {
	availableQuestions := make([]Expert.Question, 0, len(questions))
	for _, question := range questions {
		if _, excluded := excludedQuestionIDs[question.Id]; excluded {
			continue
		}
		availableQuestions = append(availableQuestions, question)
	}

	return availableQuestions
}

func buildExcludedQuestionIDs(questions []Expert.Question, answers []dbtypes.ExpertSystemAnswer) map[int]struct{} {
	excludedQuestionIDs := make(map[int]struct{}, len(answers))
	questionsByID := make(map[int]Expert.Question, len(questions))

	for _, question := range questions {
		questionsByID[question.Id] = question
	}

	for _, answer := range answers {
		excludedQuestionIDs[answer.QuestionID] = struct{}{}

		question, ok := questionsByID[answer.QuestionID]
		if !ok {
			continue
		}

		for _, blockedQuestionID := range blockedQuestionIDsForVariant(question.NoRoutes, answer.VariantIndex) {
			excludedQuestionIDs[blockedQuestionID] = struct{}{}
		}
	}

	return excludedQuestionIDs
}

func blockedQuestionIDsForVariant(noRoutes string, variantIndex int) []int {
	if strings.TrimSpace(noRoutes) == "" || variantIndex < 0 {
		return nil
	}

	variantRules := strings.Split(noRoutes, ";")
	if variantIndex >= len(variantRules) {
		return nil
	}

	rule := strings.TrimSpace(variantRules[variantIndex])
	rule = strings.Trim(rule, "\"")
	if rule == "" {
		return nil
	}

	parts := strings.Split(rule, ",")
	blockedIDs := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(strings.Trim(part, "\""))
		if part == "" {
			continue
		}

		blockedID, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		blockedIDs = append(blockedIDs, blockedID)
	}

	return blockedIDs
}

func findQuestionByID(questions []Expert.Question, questionID int) (Expert.Question, bool) {
	for _, question := range questions {
		if question.Id == questionID {
			return question, true
		}
	}
	return Expert.Question{}, false
}

func questionPosition(questions []Expert.Question, questionID int) int {
	for i, question := range questions {
		if question.Id == questionID {
			return i + 1
		}
	}
	return 0
}
