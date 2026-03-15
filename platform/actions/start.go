package actions

import (
	"TgBotUltimate/types/Action"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Start(action Action.Action) error {
	buttons := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Пройти перечень уточняющих вопросов").WithCallbackData("expert_system"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Установить параметры вручную").WithCallbackData("parameters"),
		),
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Начать подбор квартиры").WithCallbackData(ExpertSelectFlatsPrefix),
		),
	)

	_, err := action.Bot.SendMessage(
		action.ReqCtx,
		tu.Message(
			tu.ID(action.Update.Message.Chat.ID),
			"Добро пожаловать в умный телеграм бот!\n"+
				"Здесь вы сможете подобрать нужную вам квартиру, просто общаясь с нашим ботом!\n"+
				"Просто введите запрос\n\n"+
				"p.s.: Пройдите перечень уточняющих вопросов для более точных результатов",
		).WithReplyMarkup(buttons),
	)
	if err != nil {
		return err
	}
	return nil
}
