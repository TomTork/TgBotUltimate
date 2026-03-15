package actions

import (
	"TgBotUltimate/database/users"
	"TgBotUltimate/types/Action"
	dbtypes "TgBotUltimate/types/Database"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
	"strings"
	"sync"
)

const (
	ManualParametersStartPrefix = "parameters"
	ManualParameterSetPrefix    = "manual_param_set"
	ManualParameterPrevPrefix   = "manual_param_prev"
	ManualParameterNextPrefix   = "manual_param_next"
	ManualParameterFinishPrefix = "manual_param_finish"
)

type manualParameterDefinition struct {
	Key         string
	Question    string
	IsSelect    bool
	IsInteger   bool
	RangeColumn string
	GetValue    func(user *dbtypes.User) string
	SetValue    func(system *dbtypes.ExpertSystem, value string)
}

type manualParameterState struct {
	CurrentStep int
}

var manualParameterStates sync.Map

var manualParameterDefinitions = []manualParameterDefinition{
	{
		Key:      "ex_project_name",
		Question: "В каком проекте хотите начать подбор квартиры?",
		IsSelect: true,
		GetValue: func(user *dbtypes.User) string { return derefStringPointer(user.ExProjectName) },
		SetValue: func(system *dbtypes.ExpertSystem, value string) { system.ExProjectName = stringPtr(value) },
	},
	{
		Key:      "ex_building_liter",
		Question: "Какой литер вас интересует?",
		IsSelect: true,
		GetValue: func(user *dbtypes.User) string { return derefStringPointer(user.ExBuildingLiter) },
		SetValue: func(system *dbtypes.ExpertSystem, value string) { system.ExBuildingLiter = stringPtr(value) },
	},
	{
		Key:         "ex_floor_min",
		Question:    "Введите минимальный этаж.",
		IsInteger:   true,
		RangeColumn: "f.floor",
		GetValue:    func(user *dbtypes.User) string { return derefStringPointer(user.ExFloorMin) },
		SetValue:    func(system *dbtypes.ExpertSystem, value string) { system.ExFloorMin = stringPtr(value) },
	},
	{
		Key:         "ex_floor_max",
		Question:    "Введите максимальный этаж.",
		IsInteger:   true,
		RangeColumn: "f.floor",
		GetValue:    func(user *dbtypes.User) string { return derefStringPointer(user.ExFloorMax) },
		SetValue:    func(system *dbtypes.ExpertSystem, value string) { system.ExFloorMax = stringPtr(value) },
	},
	{
		Key:         "ex_rooms_amount_min",
		Question:    "Введите минимальное количество комнат.",
		IsInteger:   true,
		RangeColumn: "f.rooms_amount",
		GetValue:    func(user *dbtypes.User) string { return derefStringPointer(user.ExRoomsAmountMin) },
		SetValue:    func(system *dbtypes.ExpertSystem, value string) { system.ExRoomsAmountMin = stringPtr(value) },
	},
	{
		Key:         "ex_rooms_amount_max",
		Question:    "Введите максимальное количество комнат.",
		IsInteger:   true,
		RangeColumn: "f.rooms_amount",
		GetValue:    func(user *dbtypes.User) string { return derefStringPointer(user.ExRoomsAmountMax) },
		SetValue:    func(system *dbtypes.ExpertSystem, value string) { system.ExRoomsAmountMax = stringPtr(value) },
	},
	{
		Key:         "ex_square_min",
		Question:    "Введите минимальную площадь.",
		IsInteger:   true,
		RangeColumn: "f.total_square",
		GetValue:    func(user *dbtypes.User) string { return derefStringPointer(user.ExSquareMin) },
		SetValue:    func(system *dbtypes.ExpertSystem, value string) { system.ExSquareMin = stringPtr(value) },
	},
	{
		Key:         "ex_square_max",
		Question:    "Введите максимальную площадь.",
		IsInteger:   true,
		RangeColumn: "f.total_square",
		GetValue:    func(user *dbtypes.User) string { return derefStringPointer(user.ExSquareMax) },
		SetValue:    func(system *dbtypes.ExpertSystem, value string) { system.ExSquareMax = stringPtr(value) },
	},
	{
		Key:         "ex_cost_min",
		Question:    "Введите минимальную стоимость.",
		IsInteger:   true,
		RangeColumn: "f.cost",
		GetValue:    func(user *dbtypes.User) string { return derefStringPointer(user.ExCostMin) },
		SetValue:    func(system *dbtypes.ExpertSystem, value string) { system.ExCostMin = stringPtr(value) },
	},
	{
		Key:         "ex_cost_max",
		Question:    "Введите максимальную стоимость.",
		IsInteger:   true,
		RangeColumn: "f.cost",
		GetValue:    func(user *dbtypes.User) string { return derefStringPointer(user.ExCostMax) },
		SetValue:    func(system *dbtypes.ExpertSystem, value string) { system.ExCostMax = stringPtr(value) },
	},
}

func ManualParameters(action Action.Action) error {
	if action.Update.CallbackQuery == nil {
		return nil
	}

	callback := action.Update.CallbackQuery
	if err := answerCallback(action, callback.ID); err != nil {
		return err
	}
	if err := ensureCallbackUser(action); err != nil {
		return err
	}

	switch data := callback.Data; {
	case data == ManualParametersStartPrefix:
		return startManualParameters(action)
	case data == ManualParameterPrevPrefix:
		return moveManualParameterStep(action, -1)
	case data == ManualParameterNextPrefix:
		return moveManualParameterStep(action, 1)
	case data == ManualParameterFinishPrefix:
		return finishManualParameters(action)
	case strings.HasPrefix(data, ManualParameterSetPrefix+":"):
		return setManualParameterOption(action, data)
	default:
		return nil
	}
}

func HandleManualParameterMessage(action Action.Action) (bool, error) {
	if action.Update.Message == nil || action.Update.Message.From == nil {
		return false, nil
	}

	state := getManualParameterState(action.Update.Message.From.ID)
	if state == nil {
		return false, nil
	}
	if state.CurrentStep < 0 || state.CurrentStep >= len(manualParameterDefinitions) {
		return true, finishManualParametersByMessage(action)
	}

	definition := manualParameterDefinitions[state.CurrentStep]
	if definition.IsSelect {
		return true, sendText(action, "Для этого параметра выберите один из допустимых вариантов кнопкой.")
	}

	user, err := users.GetUserById(action.ReqCtx, action.Database, action.Update.Message.From.ID)
	if err != nil {
		return true, err
	}
	if user == nil {
		return true, nil
	}

	value := strings.TrimSpace(action.Update.Message.Text)
	if definition.IsInteger && !isStrictInteger(value) {
		return true, sendText(action, "Нужно ввести целое число без точек, запятых и других символов.")
	}

	if err := validateManualNumericValue(action, user, definition, value); err != nil {
		return true, sendText(action, err.Error())
	}

	updated := copyExpertSystem(user.ExpertSystem)
	definition.SetValue(&updated, value)
	clearManualParametersAfterStep(&updated, state.CurrentStep)

	if err := users.ReplaceExpertSystemFields(action.ReqCtx, action.Database, action.Update.Message.From.ID, updated); err != nil {
		return true, err
	}

	nextStep := state.CurrentStep + 1
	if nextStep >= len(manualParameterDefinitions) {
		return true, finishManualParametersByMessage(action)
	}

	setManualParameterState(action.Update.Message.From.ID, nextStep)

	return true, sendManualParameterStepMessage(action, nextStep)
}

func startManualParameters(action Action.Action) error {
	callback := action.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	setManualParameterState(callback.From.ID, 0)

	return renderManualParameterStep(action, 0, true)
}

func moveManualParameterStep(action Action.Action, direction int) error {
	callback := action.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	state := getManualParameterState(callback.From.ID)
	if state == nil {
		state = &manualParameterState{CurrentStep: 0}
	}

	nextStep := state.CurrentStep + direction
	if nextStep < 0 {
		nextStep = 0
	}
	if nextStep >= len(manualParameterDefinitions) {
		return finishManualParameters(action)
	}

	setManualParameterState(callback.From.ID, nextStep)

	return renderManualParameterStep(action, nextStep, true)
}

func setManualParameterOption(action Action.Action, data string) error {
	callback := action.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	parts := strings.Split(data, ":")
	if len(parts) != 3 {
		return nil
	}

	step, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}
	optionIndex, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}
	if step < 0 || step >= len(manualParameterDefinitions) {
		return nil
	}

	definition := manualParameterDefinitions[step]
	if !definition.IsSelect {
		return nil
	}

	user, err := users.GetUserById(action.ReqCtx, action.Database, callback.From.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	options, err := manualParameterOptions(action, user, definition)
	if err != nil {
		return err
	}
	if optionIndex < 0 || optionIndex >= len(options) {
		return nil
	}

	updated := copyExpertSystem(user.ExpertSystem)
	definition.SetValue(&updated, options[optionIndex])
	clearManualParametersAfterStep(&updated, step)

	if err := users.ReplaceExpertSystemFields(action.ReqCtx, action.Database, callback.From.ID, updated); err != nil {
		return err
	}
	setManualParameterState(callback.From.ID, step)

	return renderManualParameterStep(action, step, true)
}

func finishManualParameters(action Action.Action) error {
	callback := action.Update.CallbackQuery
	if callback == nil {
		return nil
	}

	deleteManualParameterState(callback.From.ID)

	_, _ = action.Bot.EditMessageText(action.ReqCtx, &telego.EditMessageTextParams{
		ChatID:    tu.ID(callback.Message.GetChat().ID),
		MessageID: callback.Message.GetMessageID(),
		Text:      "Ручная настройка параметров завершена.",
	})

	_, err := action.Bot.SendMessage(action.ReqCtx, tu.Message(
		tu.ID(callback.Message.GetChat().ID),
		"Ручная настройка параметров завершена.",
	).WithReplyMarkup(expertFinishKeyboard()))
	return err
}

func finishManualParametersByMessage(action Action.Action) error {
	if action.Update.Message == nil {
		return nil
	}

	deleteManualParameterState(action.Update.Message.From.ID)

	_, err := action.Bot.SendMessage(action.ReqCtx, tu.Message(
		tu.ID(action.Update.Message.Chat.ID),
		"Ручная настройка параметров завершена.",
	).WithReplyMarkup(expertFinishKeyboard()))
	return err
}

func renderManualParameterStep(action Action.Action, step int, editCurrentMessage bool) error {
	callback := action.Update.CallbackQuery
	if callback == nil || callback.Message == nil {
		return nil
	}
	if step < 0 || step >= len(manualParameterDefinitions) {
		return finishManualParameters(action)
	}

	user, err := users.GetUserById(action.ReqCtx, action.Database, callback.From.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	text, keyboard, err := buildManualParameterView(action, user, step)
	if err != nil {
		return err
	}

	if editCurrentMessage {
		_, err = action.Bot.EditMessageText(action.ReqCtx, &telego.EditMessageTextParams{
			ChatID:      tu.ID(callback.Message.GetChat().ID),
			MessageID:   callback.Message.GetMessageID(),
			Text:        text,
			ReplyMarkup: keyboard,
		})
		return err
	}

	_, err = action.Bot.SendMessage(action.ReqCtx, tu.Message(
		tu.ID(callback.Message.GetChat().ID),
		text,
	).WithReplyMarkup(keyboard))
	return err
}

func sendManualParameterStepMessage(action Action.Action, step int) error {
	if action.Update.Message == nil {
		return nil
	}
	if step < 0 || step >= len(manualParameterDefinitions) {
		return finishManualParametersByMessage(action)
	}

	user, err := users.GetUserById(action.ReqCtx, action.Database, action.Update.Message.From.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	text, keyboard, err := buildManualParameterView(action, user, step)
	if err != nil {
		return err
	}

	_, err = action.Bot.SendMessage(action.ReqCtx, tu.Message(
		tu.ID(action.Update.Message.Chat.ID),
		text,
	).WithReplyMarkup(keyboard))
	return err
}

func buildManualParameterView(action Action.Action, user *dbtypes.User, step int) (string, *telego.InlineKeyboardMarkup, error) {
	definition := manualParameterDefinitions[step]
	currentValue := definition.GetValue(user)

	textParts := []string{
		fmt.Sprintf("Параметр %d из %d", step+1, len(manualParameterDefinitions)),
		"",
		definition.Question,
	}
	if currentValue != "" {
		textParts = append(textParts, fmt.Sprintf("Текущее значение: %s", currentValue))
	}

	buttonRows := make([][]telego.InlineKeyboardButton, 0)

	if definition.IsSelect {
		options, err := manualParameterOptions(action, user, definition)
		if err != nil {
			return "", nil, err
		}
		if len(options) == 0 {
			textParts = append(textParts, "Подходящих вариантов нет. Измените предыдущие параметры или выполните сброс.")
		}
		for index, option := range options {
			label := option
			if option == currentValue {
				label = "• " + option
			}
			buttonRows = append(buttonRows, tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(label).
					WithCallbackData(fmt.Sprintf("%s:%d:%d", ManualParameterSetPrefix, step, index)),
			))
		}
	} else {
		textParts = append(textParts, "Введите значение сообщением.")
		if definition.IsInteger {
			minValue, maxValue, hasRange, err := manualParameterRange(action, user, definition)
			if err != nil {
				return "", nil, err
			}
			if hasRange {
				textParts = append(textParts, fmt.Sprintf("Допустимый диапазон: %d-%d", minValue, maxValue))
			}
		}
	}

	buttonRows = append(buttonRows, tu.InlineKeyboardRow(
		tu.InlineKeyboardButton("Сбросить").WithCallbackData(ExpertResetPrefix),
	))

	navigationRow := make([]telego.InlineKeyboardButton, 0, 3)
	if step > 0 {
		navigationRow = append(navigationRow, tu.InlineKeyboardButton("◀").WithCallbackData(ManualParameterPrevPrefix))
	}
	navigationRow = append(navigationRow, tu.InlineKeyboardButton("Завершить").WithCallbackData(ManualParameterFinishPrefix))
	if step < len(manualParameterDefinitions)-1 {
		navigationRow = append(navigationRow, tu.InlineKeyboardButton("▶").WithCallbackData(ManualParameterNextPrefix))
	}
	buttonRows = append(buttonRows, navigationRow)

	return strings.Join(textParts, "\n"), tu.InlineKeyboard(buttonRows...), nil
}

func manualParameterOptions(action Action.Action, user *dbtypes.User, definition manualParameterDefinition) ([]string, error) {
	selectExpr := ""
	switch definition.Key {
	case "ex_project_name":
		selectExpr = `SELECT DISTINCT p.name `
	case "ex_building_liter":
		selectExpr = `SELECT DISTINCT b.liter `
	default:
		return nil, nil
	}

	whereClause, args := buildManualParameterFilters(user, map[string]struct{}{
		definition.Key: {},
	})
	orderBy := " ORDER BY p.name"
	if definition.Key == "ex_building_liter" {
		whereClause += " AND b.liter <> ''"
		orderBy = " ORDER BY b.liter"
	}

	rows, err := action.Database.Query(action.ReqCtx, selectExpr+manualParameterBaseQuery()+whereClause+orderBy, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := make([]string, 0)
	for rows.Next() {
		var option string
		if err := rows.Scan(&option); err != nil {
			return nil, err
		}
		option = strings.TrimSpace(option)
		if option == "" {
			continue
		}
		options = append(options, option)
	}

	return options, rows.Err()
}

func manualParameterRange(action Action.Action, user *dbtypes.User, definition manualParameterDefinition) (int, int, bool, error) {
	if definition.RangeColumn == "" {
		return 0, 0, false, nil
	}

	rangeExprMin := fmt.Sprintf("MIN(%s)::numeric", definition.RangeColumn)
	rangeExprMax := fmt.Sprintf("MAX(%s)::numeric", definition.RangeColumn)
	if definition.RangeColumn == "f.total_square" || definition.RangeColumn == "f.cost" {
		rangeExprMin = fmt.Sprintf("CEIL(MIN(%s))::integer", definition.RangeColumn)
		rangeExprMax = fmt.Sprintf("FLOOR(MAX(%s))::integer", definition.RangeColumn)
	} else {
		rangeExprMin = fmt.Sprintf("MIN(%s)::integer", definition.RangeColumn)
		rangeExprMax = fmt.Sprintf("MAX(%s)::integer", definition.RangeColumn)
	}

	skip := map[string]struct{}{
		definition.Key: {},
	}
	if pairKey := manualParameterPairKey(definition.Key); pairKey != "" {
		skip[pairKey] = struct{}{}
	}
	whereClause, args := buildManualParameterFilters(user, skip)

	query := fmt.Sprintf(`SELECT %s, %s `, rangeExprMin, rangeExprMax) + manualParameterBaseQuery() + whereClause

	var minValue, maxValue *int
	if err := action.Database.QueryRow(action.ReqCtx, query, args...).Scan(&minValue, &maxValue); err != nil {
		return 0, 0, false, err
	}
	if minValue == nil || maxValue == nil {
		return 0, 0, false, nil
	}

	return *minValue, *maxValue, true, nil
}

func validateManualNumericValue(action Action.Action, user *dbtypes.User, definition manualParameterDefinition, value string) error {
	if !definition.IsInteger {
		return nil
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("Нужно ввести целое число без точек, запятых и других символов.")
	}

	minValue, maxValue, hasRange, err := manualParameterRange(action, user, definition)
	if err != nil {
		return err
	}
	if hasRange && (intValue < minValue || intValue > maxValue) {
		return fmt.Errorf("Введите значение в допустимом диапазоне: %d-%d.", minValue, maxValue)
	}

	if pairKey := manualParameterPairKey(definition.Key); pairKey != "" {
		pairValue := manualParameterValueByKey(user, pairKey)
		if pairValue != "" {
			pairInt, err := strconv.Atoi(pairValue)
			if err == nil {
				if strings.HasSuffix(definition.Key, "_min") && intValue > pairInt {
					return fmt.Errorf("Минимальное значение не может быть больше максимального (%d).", pairInt)
				}
				if strings.HasSuffix(definition.Key, "_max") && intValue < pairInt {
					return fmt.Errorf("Максимальное значение не может быть меньше минимального (%d).", pairInt)
				}
			}
		}
	}

	return nil
}

func buildManualParameterFilters(user *dbtypes.User, skip map[string]struct{}) (string, []interface{}) {
	conditions := []string{" WHERE f.status = 0"}
	args := make([]interface{}, 0)

	for _, definition := range manualParameterDefinitions {
		if _, skipped := skip[definition.Key]; skipped {
			continue
		}

		value := definition.GetValue(user)
		if value == "" {
			continue
		}

		args = append(args, value)
		placeholder := fmt.Sprintf("$%d", len(args))

		switch definition.Key {
		case "ex_project_name":
			conditions = append(conditions, " AND p.name = "+placeholder)
		case "ex_building_liter":
			conditions = append(conditions, " AND b.liter = "+placeholder)
		case "ex_floor_min":
			conditions = append(conditions, " AND f.floor >= "+placeholder)
		case "ex_floor_max":
			conditions = append(conditions, " AND f.floor <= "+placeholder)
		case "ex_rooms_amount_min":
			conditions = append(conditions, " AND f.rooms_amount >= "+placeholder)
		case "ex_rooms_amount_max":
			conditions = append(conditions, " AND f.rooms_amount <= "+placeholder)
		case "ex_square_min":
			conditions = append(conditions, " AND f.total_square >= "+placeholder)
		case "ex_square_max":
			conditions = append(conditions, " AND f.total_square <= "+placeholder)
		case "ex_cost_min":
			conditions = append(conditions, " AND f.cost >= "+placeholder)
		case "ex_cost_max":
			conditions = append(conditions, " AND f.cost <= "+placeholder)
		}
	}

	return strings.Join(conditions, ""), args
}

func manualParameterBaseQuery() string {
	return `FROM flats f
LEFT JOIN buildings b ON b.code = f.building_code
LEFT JOIN projects p ON p.code = b.project_code`
}

func manualParameterPairKey(key string) string {
	switch key {
	case "ex_floor_min":
		return "ex_floor_max"
	case "ex_floor_max":
		return "ex_floor_min"
	case "ex_rooms_amount_min":
		return "ex_rooms_amount_max"
	case "ex_rooms_amount_max":
		return "ex_rooms_amount_min"
	case "ex_square_min":
		return "ex_square_max"
	case "ex_square_max":
		return "ex_square_min"
	case "ex_cost_min":
		return "ex_cost_max"
	case "ex_cost_max":
		return "ex_cost_min"
	default:
		return ""
	}
}

func manualParameterValueByKey(user *dbtypes.User, key string) string {
	for _, definition := range manualParameterDefinitions {
		if definition.Key == key {
			return definition.GetValue(user)
		}
	}
	return ""
}

func clearManualParametersAfterStep(system *dbtypes.ExpertSystem, step int) {
	for index := step + 1; index < len(manualParameterDefinitions); index++ {
		manualParameterDefinitions[index].SetValue(system, "")
	}
}

func copyExpertSystem(system dbtypes.ExpertSystem) dbtypes.ExpertSystem {
	return dbtypes.ExpertSystem{
		ExProjectName:    stringPointerCopy(system.ExProjectName),
		ExBuildingLiter:  stringPointerCopy(system.ExBuildingLiter),
		ExFloorMin:       stringPointerCopy(system.ExFloorMin),
		ExFloorMax:       stringPointerCopy(system.ExFloorMax),
		ExRoomsAmountMin: stringPointerCopy(system.ExRoomsAmountMin),
		ExRoomsAmountMax: stringPointerCopy(system.ExRoomsAmountMax),
		ExSquareMin:      stringPointerCopy(system.ExSquareMin),
		ExSquareMax:      stringPointerCopy(system.ExSquareMax),
		ExCostMin:        stringPointerCopy(system.ExCostMin),
		ExCostMax:        stringPointerCopy(system.ExCostMax),
	}
}

func stringPointerCopy(value *string) *string {
	if value == nil {
		return stringPtr("")
	}
	return stringPtr(*value)
}

func getManualParameterState(userID int64) *manualParameterState {
	value, ok := manualParameterStates.Load(userID)
	if !ok {
		return nil
	}

	state, ok := value.(manualParameterState)
	if !ok {
		return nil
	}

	return &state
}

func setManualParameterState(userID int64, step int) {
	manualParameterStates.Store(userID, manualParameterState{CurrentStep: step})
}

func deleteManualParameterState(userID int64) {
	manualParameterStates.Delete(userID)
}

func derefStringPointer(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func isStrictInteger(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
