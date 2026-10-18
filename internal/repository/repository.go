package repository

import (
	"cinema-booking/internal/domain"
	"context"
	"time"
)

type ShowtimeRepository interface {
	List(ctx context.Context) ([]domain.Showtime, error)
	Get(ctx context.Context, id int64) (*domain.Showtime, error)

	SeatsByHall(ctx context.Context, hallID int64) ([]domain.Seat, error)
	BookedSeatIDs(ctx context.Context, showtimeID int64) (map[int64]bool, error)
}

type SeatHoldStore interface {
	Hold(ctx context.Context, showtimeID int64, seatIDs []int64, userID string, ttl time.Duration) error
	Release(ctx context.Context, showtimeID int64, seatIDs []int64, userID string) error
	Holds(ctx context.Context, showtimeID int64, seatIDs []int64) (map[int64]string, error)
}
