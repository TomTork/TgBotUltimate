package neuro

import (
	"TgBotUltimate/server/routes/helper"
	"TgBotUltimate/types/Ollama"
	"context"
	"fmt"
	"log"
	"os"
)

func Ask(ctx context.Context, prompt string) string {
	return ask(ctx, prompt, "qwen3-vl:4b")
}

func Parameters(ctx context.Context, prompt string) string {
	log.Println(fmt.Sprintf(
		"Ты должен вычленить параметры подбора недвижимости на основе информации: %s\n"+
			"и заполнить их по форме ниже:\n"+
			"project_name: <UNK>"+
			"building_liter: <UNK>"+
			"floor_min: <UNK>"+
			"floor_max: <UNK>"+
			"rooms_amount_min: <UNK>"+
			"rooms_amount_max: <UNK>"+
			"square_min: <UNK>"+
			"square_max: <UNK>"+
			"Вместо <UNK> подставь значения на основе запроса пользователя. "+
			"Выведи ТОЛЬКО форму.",
		prompt,
	))
	return ask(ctx,
		fmt.Sprintf(
			"Ты должен вычленить параметры подбора недвижимости на основе информации: %s\n"+
				"и заполнить их по форме ниже:\n"+
				"project_name: <UNK>"+
				"building_liter: <UNK>"+
				"floor_min: <UNK>"+
				"floor_max: <UNK>"+
				"rooms_amount_min: <UNK>"+
				"rooms_amount_max: <UNK>"+
				"square_min: <UNK>"+
				"square_max: <UNK>"+
				"Вместо <UNK> подставь значения на основе запроса пользователя. "+
				"Выведи ТОЛЬКО форму.",
			prompt,
		),
		"qwen3-embedding:0.6b",
	)
}

func ask(ctx context.Context, prompt string, model string) string {
	var response Ollama.Response
	err := helper.Post(
		ctx,
		fmt.Sprintf("http://localhost:%s/api/generate", os.Getenv("NEURO_PORT")),
		nil,
		Ollama.Request{
			Model:  model,
			Prompt: prompt,
			Stream: false,
			Options: map[string]interface{}{
				"temperature": 0.7,
			},
		},
		&response,
	)
	if err != nil {
		return err.Error()
	}
	log.Println(response)
	return response.Response
}
