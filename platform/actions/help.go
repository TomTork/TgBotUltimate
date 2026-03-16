package actions

import (
	"TgBotUltimate/types/Action"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Help(action Action.Action) error {
	_, err := action.Bot.SendMessage(action.ReqCtx, tu.Message(tu.ID(action.Update.Message.Chat.ID), "Бот работает по принципу объединения экспертной системы и запросов пользователя.\n\n"+
		"То есть вы можете ответить на вопросы экспертной системы и они приоритетно будут влиять на выдачу квартир.\n"+
		"При этом у вас есть возможность вводить запросы просто с клавиатуры, обученная нейросеть возьмёт из вашего запроса всё необходимое.\n\n"+
		"Давайте же приступим к подбору квартиры мечты!\n\n"+"Доступные команды:\n"+
		"/start — Приветственное окно\n"+
		"/help — Это сообщение\n"+
		"/questions — Пройти перечень уточняющих вопросов\n"+
		"/reload — Сбросить параметры подбора\n"+
		"/flats — Начать подбор квартир\n"+
		"/favorites — Показать избранные планировки"))
	if err != nil {
		return err
	}
	return nil
}
