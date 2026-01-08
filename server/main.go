package server

import (
	"TgBotUltimate/server/routes"
	"fmt"
	"net/http"
	"os"
)

func RunHTTP() error {
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), routes.NewRouter())
	if err != nil {
		return err
	}
	return nil
}
