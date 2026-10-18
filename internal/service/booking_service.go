package service

import (
	"context"
	"fmt"

	"cinema-booking/internal/domain"
	"cinema-booking/internal/repository"
)

type BookingService struct {
	showtimes repository.ShowtimeRepository
}

func NewBookingService(showtimes repository.ShowtimeRepository) *BookingService {
	return &BookingService{showtimes: showtimes}
}

func (s *BookingService) SeatMap(ctx context.Context, showtimeID int64) ([]domain.SeatView, error) {
	st, err := s.showtimes.Get(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	seats, err := s.showtimes.SeatsByHall(ctx, st.HallID)
	if err != nil {
		return nil, err
	}

	booked, err := s.showtimes.BookedSeatIDs(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	out := make([]domain.SeatView, 0, len(seats))
	for _, seat := range seats {
		status := domain.SeatAvailable
		if booked[seat.ID] {
			status = domain.SeatBooked
		}
		out = append(out, domain.SeatView{
			ID:     seat.ID,
			Label:  fmt.Sprintf("%s%d", seat.Row, seat.Number),
			Status: status,
		})
	}
	return out, nil
}
