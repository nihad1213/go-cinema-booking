package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"cinema-booking/internal/domain"
)

type BookingRepo struct{ db *sql.DB }

func NewBookingRepo(db *sql.DB) *BookingRepo { return &BookingRepo{db: db} }

func (r *BookingRepo) CreateConfirmed(ctx context.Context, b *domain.Booking) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	res, err := tx.ExecContext(ctx,
		`INSERT INTO bookings (user_id, showtime_id, status, total_cents, created_at)
		 VALUES (?, ?, 'confirmed', ?, ?)`,
		b.UserID, b.ShowtimeID, b.TotalCents, now.Unix())
	if err != nil {
		return err
	}
	bookingID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, seatID := range b.SeatIDs {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO booking_seats (booking_id, showtime_id, seat_id) VALUES (?, ?, ?)`,
			bookingID, b.ShowtimeID, seatID)
		if err != nil {
			if isUniqueViolation(err) {
				return domain.ErrSeatUnavailable
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	b.ID, b.Status, b.CreatedAt = bookingID, "confirmed", now
	return nil
}

func (r *BookingRepo) Get(ctx context.Context, id int64) (*domain.Booking, error) {
	var b domain.Booking
	var ts int64
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, showtime_id, status, total_cents, created_at
		 FROM bookings WHERE id = ?`, id).
		Scan(&b.ID, &b.UserID, &b.ShowtimeID, &b.Status, &b.TotalCents, &ts)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	b.CreatedAt = time.Unix(ts, 0).UTC()

	rows, err := r.db.QueryContext(ctx,
		`SELECT seat_id FROM booking_seats WHERE booking_id = ? ORDER BY seat_id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	b.SeatIDs = []int64{}
	for rows.Next() {
		var seatID int64
		if err := rows.Scan(&seatID); err != nil {
			return nil, err
		}
		b.SeatIDs = append(b.SeatIDs, seatID)
	}
	return &b, rows.Err()
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
