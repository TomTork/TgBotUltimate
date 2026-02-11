package health

import (
	"TgBotUltimate/database"
	"TgBotUltimate/database/data"
	"github.com/go-chi/chi/v5"
	"github.com/grbit/go-json"
	"net/http"
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

	return r
}
