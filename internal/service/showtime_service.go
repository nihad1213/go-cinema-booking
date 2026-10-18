package service

import (
	"context"

	"cinema-booking/internal/domain"
	"cinema-booking/internal/repository"
)

type ShowtimeService struct {
	showtimes repository.ShowtimeRepository
}

func NewShowtimeService(showtimes repository.ShowtimeRepository) *ShowtimeService {
	return &ShowtimeService{showtimes: showtimes}
}

func (s *ShowtimeService) List(ctx context.Context) ([]domain.Showtime, error) {
	return s.showtimes.List(ctx)
}

func (s *ShowtimeService) Get(ctx context.Context, id int64) (*domain.Showtime, error) {
	return s.showtimes.Get(ctx, id)
}
