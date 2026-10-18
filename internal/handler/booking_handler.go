package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cinema-booking/internal/domain"
	"cinema-booking/internal/service"
)

type BookingHandler struct{ svc *service.BookingService }

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

type seatsRequest struct {
	SeatIDs []int64 `json:"seat_ids"`
}

func decodeSeats(r *http.Request) ([]int64, error) {
	var req seatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("%w: body must be JSON like {\"seat_ids\":[1,2]}", domain.ErrInvalidInput)
	}
	return req.SeatIDs, nil
}

func (h *BookingHandler) SeatMap(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	out, err := h.svc.SeatMap(r.Context(), id, currentUser(r))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *BookingHandler) Hold(w http.ResponseWriter, r *http.Request) {
	user, err := requireUser(r)
	if err != nil {
		writeError(w, err)
		return
	}
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	seats, err := decodeSeats(r)
	if err != nil {
		writeError(w, err)
		return
	}
	res, err := h.svc.HoldSeats(r.Context(), id, user, seats)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *BookingHandler) Release(w http.ResponseWriter, r *http.Request) {
	user, err := requireUser(r)
	if err != nil {
		writeError(w, err)
		return
	}
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	seats, err := decodeSeats(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if err := h.svc.ReleaseSeats(r.Context(), id, user, seats); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *BookingHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	user, err := requireUser(r)
	if err != nil {
		writeError(w, err)
		return
	}
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	seats, err := decodeSeats(r)
	if err != nil {
		writeError(w, err)
		return
	}
	booking, err := h.svc.ConfirmBooking(r.Context(), id, user, seats)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, booking)
}

func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	booking, err := h.svc.GetBooking(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, booking)
}
