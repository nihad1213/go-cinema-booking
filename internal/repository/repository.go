package repository

import (
	"context"

	"cinema-booking/internal/domain"
)

type ShowtimeRepository interface {
	List(ctx context.Context) ([]domain.Showtime, error)
	Get(ctx context.Context, id int64) (*domain.Showtime, error)
}
