package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cinema-booking/internal/domain"
)

type ShowtimeRepo struct{ db *sql.DB }

func NewShowtimeRepo(db *sql.DB) *ShowtimeRepo { return &ShowtimeRepo{db: db} }

const showtimeSelect = `
SELECT s.id, s.movie_id, m.title, s.hall_id, s.starts_at, s.price_cents
FROM showtimes s
JOIN movies m ON m.id = s.movie_id`

type scanner interface{ Scan(dest ...any) error }

func scanShowtime(sc scanner) (*domain.Showtime, error) {
	var st domain.Showtime
	var ts int64
	if err := sc.Scan(&st.ID, &st.MovieID, &st.MovieTitle, &st.HallID, &ts, &st.PriceCents); err != nil {
		return nil, err
	}
	st.StartsAt = time.Unix(ts, 0).UTC()
	return &st, nil
}

func (r *ShowtimeRepo) List(ctx context.Context) ([]domain.Showtime, error) {
	rows, err := r.db.QueryContext(ctx, showtimeSelect+` ORDER BY s.starts_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Showtime{}
	for rows.Next() {
		st, err := scanShowtime(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *st)
	}
	return out, rows.Err()
}

func (r *ShowtimeRepo) Get(ctx context.Context, id int64) (*domain.Showtime, error) {
	st, err := scanShowtime(r.db.QueryRowContext(ctx, showtimeSelect+` WHERE s.id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return st, err
}
