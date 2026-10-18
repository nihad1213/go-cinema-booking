package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"cinema-booking/internal/domain"
	"cinema-booking/internal/repository"
)

type BookingService struct {
	showtimes repository.ShowtimeRepository
	holds     repository.SeatHoldStore
	holdTTL   time.Duration
	maxSeats  int
}

func NewBookingService(
	showtimes repository.ShowtimeRepository,
	holds repository.SeatHoldStore,
	holdTTL time.Duration,
	maxSeats int,
) *BookingService {
	return &BookingService{
		showtimes: showtimes,
		holds:     holds,
		holdTTL:   holdTTL,
		maxSeats:  maxSeats,
	}
}

func (s *BookingService) SeatMap(ctx context.Context, showtimeID int64, userID string) ([]domain.SeatView, error) {
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

	ids := make([]int64, len(seats))
	for i, seat := range seats {
		ids[i] = seat.ID
	}
	held, err := s.holds.Holds(ctx, showtimeID, ids)
	if err != nil {
		return nil, err
	}

	out := make([]domain.SeatView, 0, len(seats))
	for _, seat := range seats {
		v := domain.SeatView{
			ID:     seat.ID,
			Label:  fmt.Sprintf("%s%d", seat.Row, seat.Number),
			Status: domain.SeatAvailable,
		}
		if booked[seat.ID] {
			v.Status = domain.SeatBooked
		} else if owner, ok := held[seat.ID]; ok {
			v.Status = domain.SeatPending
			v.Mine = owner == userID
		}
		out = append(out, v)
	}
	return out, nil
}

func (s *BookingService) HoldSeats(ctx context.Context, showtimeID int64, userID string, seatIDs []int64) (*domain.HoldResult, error) {
	seatIDs, err := s.normalize(seatIDs)
	if err != nil {
		return nil, err
	}
	if err := s.validateSeats(ctx, showtimeID, seatIDs); err != nil {
		return nil, err
	}
	if err := s.holds.Hold(ctx, showtimeID, seatIDs, userID, s.holdTTL); err != nil {
		return nil, err
	}
	return &domain.HoldResult{
		ShowtimeID: showtimeID,
		SeatIDs:    seatIDs,
		ExpiresAt:  time.Now().Add(s.holdTTL).UTC(),
	}, nil
}

func (s *BookingService) ReleaseSeats(ctx context.Context, showtimeID int64, userID string, seatIDs []int64) error {
	seatIDs, err := s.normalize(seatIDs)
	if err != nil {
		return err
	}
	return s.holds.Release(ctx, showtimeID, seatIDs, userID)
}

func (s *BookingService) normalize(ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("%w: seat_ids is empty", domain.ErrInvalidInput)
	}
	if len(ids) > s.maxSeats {
		return nil, fmt.Errorf("%w: max %d seats per booking", domain.ErrInvalidInput, s.maxSeats)
	}

	seen := make(map[int64]bool, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			return nil, fmt.Errorf("%w: bad seat id %d", domain.ErrInvalidInput, id)
		}
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out, nil
}

func (s *BookingService) validateSeats(ctx context.Context, showtimeID int64, seatIDs []int64) error {
	st, err := s.showtimes.Get(ctx, showtimeID)
	if err != nil {
		return err
	}
	seats, err := s.showtimes.SeatsByHall(ctx, st.HallID)
	if err != nil {
		return err
	}
	inHall := make(map[int64]bool, len(seats))
	for _, seat := range seats {
		inHall[seat.ID] = true
	}
	for _, id := range seatIDs {
		if !inHall[id] {
			return fmt.Errorf("%w: seat %d is not in this hall", domain.ErrInvalidInput, id)
		}
	}

	booked, err := s.showtimes.BookedSeatIDs(ctx, showtimeID)
	if err != nil {
		return err
	}
	for _, id := range seatIDs {
		if booked[id] {
			return domain.ErrSeatUnavailable
		}
	}
	return nil
}
