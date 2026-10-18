CREATE TABLE movies (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    title        TEXT    NOT NULL,
    duration_min INTEGER NOT NULL
);

CREATE TABLE halls (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE seats (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    hall_id   INTEGER NOT NULL REFERENCES halls(id),
    row_label TEXT    NOT NULL,
    number    INTEGER NOT NULL,
    UNIQUE (hall_id, row_label, number)
);

CREATE TABLE showtimes (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    movie_id    INTEGER NOT NULL REFERENCES movies(id),
    hall_id     INTEGER NOT NULL REFERENCES halls(id),
    starts_at   INTEGER NOT NULL,
    price_cents INTEGER NOT NULL
);

CREATE TABLE bookings (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     TEXT    NOT NULL,
    showtime_id INTEGER NOT NULL REFERENCES showtimes(id),
    status      TEXT    NOT NULL DEFAULT 'confirmed',
    total_cents INTEGER NOT NULL,
    created_at  INTEGER NOT NULL
);

CREATE TABLE booking_seats (
    booking_id  INTEGER NOT NULL REFERENCES bookings(id),
    showtime_id INTEGER NOT NULL REFERENCES showtimes(id),
    seat_id     INTEGER NOT NULL REFERENCES seats(id),
    UNIQUE (showtime_id, seat_id)
);