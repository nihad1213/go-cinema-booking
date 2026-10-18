package repository

import (
	"context"

	"cinema-booking/internal/domain"
)

type ShowtimeRepository interface {
	List(ctx context.Context) ([]domain.Showtime, error)
	Get(ctx context.Context, id int64) (*domain.Showtime, error)

	SeatsByHall(ctx context.Context, hallID int64) ([]domain.Seat, error)
	BookedSeatIDs(ctx context.Context, showtimeID int64) (map[int64]bool, error)
}
