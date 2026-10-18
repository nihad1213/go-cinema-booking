package repository

import (
	"context"
	"time"

	"cinema-booking/internal/domain"
)

type ShowtimeRepository interface {
	List(ctx context.Context) ([]domain.Showtime, error)
	Get(ctx context.Context, id int64) (*domain.Showtime, error)
	SeatsByHall(ctx context.Context, hallID int64) ([]domain.Seat, error)
	BookedSeatIDs(ctx context.Context, showtimeID int64) (map[int64]bool, error)
}

type BookingRepository interface {
	CreateConfirmed(ctx context.Context, b *domain.Booking) error
	Get(ctx context.Context, id int64) (*domain.Booking, error)
}

type SeatHoldStore interface {
	Hold(ctx context.Context, showtimeID int64, seatIDs []int64, userID string, ttl time.Duration) error
	Release(ctx context.Context, showtimeID int64, seatIDs []int64, userID string) error
	Holds(ctx context.Context, showtimeID int64, seatIDs []int64) (map[int64]string, error)
	Verify(ctx context.Context, showtimeID int64, seatIDs []int64, userID string, extend time.Duration) (bool, error)
}
