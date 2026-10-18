package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    HTTPAddr   string
    SQLitePath string
    HoldTTL    time.Duration
    MaxSeats   int
}

func Load() Config {
    return Config{
        HTTPAddr:   env("HTTP_ADDR", ":8080"),
        SQLitePath: env("SQLITE_PATH", "data/cinema.db"),
        HoldTTL:    envDuration("HOLD_TTL", 5*time.Minute),
        MaxSeats:   envInt("MAX_SEATS_PER_BOOKING", 6),
    }
}

func env(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func envInt(key string, def int) int {
    if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
        return v
    }
    return def
}

func envDuration(key string, def time.Duration) time.Duration {
    if v, err := time.ParseDuration(os.Getenv(key)); err == nil {
        return v
    }
    return def
}