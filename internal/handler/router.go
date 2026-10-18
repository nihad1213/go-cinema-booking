package handler

import "net/http"

func NewRouter(showtimes *ShowtimeHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /api/showtimes", showtimes.List)
	mux.HandleFunc("GET /api/showtimes/{id}", showtimes.Get)

	return mux
}
