package actions

import (
	"TgBotUltimate/types/Action"
	"github.com/mymmrac/telego"
	"log"
)

func SetCommands(action Action.Action) {
	commands := []telego.BotCommand{
		{Command: "start", Description: "Приветственное окно"},
		{Command: "help", Description: "Как это работает"},
		{Command: "questions", Description: "Пройти уточняющие вопросы"},
		{Command: "reload", Description: "Сбросить параметры"},
		{Command: "flats", Description: "Начать подбор квартиры"},
		{Command: "favorites", Description: "Показать избранные планировки"},
	}
	err := action.Bot.SetMyCommands(action.Ctx, &telego.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		log.Println(err)
	}
}
