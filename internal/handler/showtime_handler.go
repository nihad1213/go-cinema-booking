package handler

import (
	"net/http"

	"cinema-booking/internal/service"
)

type ShowtimeHandler struct{ svc *service.ShowtimeService }

func NewShowtimeHandler(svc *service.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{svc: svc}
}

func (h *ShowtimeHandler) List(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *ShowtimeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	out, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
