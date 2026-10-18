package handler

import "net/http"

func NewRouter(showtimes *ShowtimeHandler, bookings *BookingHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /api/showtimes", showtimes.List)
	mux.HandleFunc("GET /api/showtimes/{id}", showtimes.Get)
	mux.HandleFunc("GET /api/showtimes/{id}/seats", bookings.SeatMap)
	mux.HandleFunc("POST /api/showtimes/{id}/hold", bookings.Hold)
	mux.HandleFunc("DELETE /api/showtimes/{id}/hold", bookings.Release)
	mux.HandleFunc("POST /api/showtimes/{id}/confirm", bookings.Confirm)
	mux.HandleFunc("GET /api/bookings/{id}", bookings.GetBooking)

	mux.Handle("GET /", http.FileServer(http.Dir("web")))

	return UserID(mux)
}
