package handler

import (
	"TgBotUltimate/database"
	"TgBotUltimate/database/data"
	"TgBotUltimate/server/routes/external/core"
	"github.com/go-chi/chi/v5"
	"github.com/grbit/go-json"
	"net/http"
	"os"
)

func Handler() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		db, err := database.NewDatabase(r.Context())
		flat, err := data.GetFlatByCode(r.Context(), db, "2")
		ser, err := json.Marshal(flat)
		_, err = w.Write([]byte(ser))
		if err != nil {
			return
		}
	})

	r.Post("/sync", func(w http.ResponseWriter, r *http.Request) {
		switch os.Getenv("SYNC") {
		case "FEED":
			_, err := w.Write([]byte(core.Feed(r.Context())))
			if err != nil {
				return
			}
		case "STRAPI":
			_, err := w.Write([]byte(core.Strapi(r.Context())))
			if err != nil {
				return
			}
		default:
			_, err := w.Write([]byte("Error: Undefined cron task."))
			if err != nil {
				return
			}
		}
	})

	return r
}
