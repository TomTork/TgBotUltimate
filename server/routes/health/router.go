package health

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Handler() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			return
		}
	})

	return r
}
