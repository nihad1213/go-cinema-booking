package domain

import "time"

type SeatStatus string

const (
	SeatAvailable SeatStatus = "available"
	SeatPending   SeatStatus = "pending"
	SeatBooked    SeatStatus = "booked"
)

type SeatView struct {
	ID     int64      `json:"id"`
	Label  string     `json:"label"`
	Status SeatStatus `json:"status"`
}

type Showtime struct {
	ID         int64     `json:"id"`
	MovieID    int64     `json:"movie_id"`
	MovieTitle string    `json:"movie_title"`
	HallID     int64     `json:"hall_id"`
	StartsAt   time.Time `json:"starts_at"`
	PriceCents int64     `json:"price_cents"`
}

type Seat struct {
	ID     int64  `json:"id"`
	HallID int64  `json:"hall_id"`
	Row    string `json:"row"`
	Number int    `json:"number"`
}

type Booking struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	ShowtimeID int64     `json:"showtime_id"`
	SeatIDs    []int64   `json:"seat_ids"`
	TotalCents int64     `json:"total_cents"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
