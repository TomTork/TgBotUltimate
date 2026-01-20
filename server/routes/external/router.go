package external

import (
	"TgBotUltimate/server/routes/external/core"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

func Handler() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/sync", func(w http.ResponseWriter, r *http.Request) {
		switch os.Getenv("SYNC") {
		case "FEED":
			_, err := w.Write([]byte(core.Feed(r.Context())))
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
