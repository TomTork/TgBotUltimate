package actions

import (
	"TgBotUltimate/types/Action"
	"log"
	"strings"
)

func CallbackQuery(action Action.Action) {
	if action.Update.CallbackQuery == nil {
		return
	}
	data := action.Update.CallbackQuery.Data
	switch {
	case data == "expert_system",
		data == "expert_finish",
		data == "expert_reset",
		data == "expert_select_flats",
		strings.HasPrefix(data, "expert_answer:"),
		strings.HasPrefix(data, "expert_prev:"),
		strings.HasPrefix(data, "expert_next:"):
		if err := ExpertSystem(action); err != nil {
			log.Println("expert_system error:", err)
		}
	case data == ShowMoreFlatsPrefix:
		if err := ShowMoreFlats(action); err != nil {
			log.Println("show_more_flats error:", err)
		}
	}
}
