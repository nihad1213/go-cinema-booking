package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"cinema-booking/internal/config"
	"cinema-booking/internal/handler"
	"cinema-booking/internal/platform/redisclient"
	"cinema-booking/internal/platform/sqlitedb"
	"cinema-booking/internal/repository/rediscache"
	"cinema-booking/internal/repository/sqlite"
	"cinema-booking/internal/service"
	"cinema-booking/migrations"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	if err := os.MkdirAll(filepath.Dir(cfg.SQLitePath), 0o755); err != nil {
		log.Fatal(err)
	}
	db, err := sqlitedb.Open(cfg.SQLitePath)
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	if err := sqlitedb.Migrate(ctx, db, migrations.FS); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	rdb, err := redisclient.New(ctx, cfg.RedisAddr)
	if err != nil {
		log.Fatalf("redis: %v (is it running? try: docker compose up -d redis)", err)
	}
	defer rdb.Close()

	showtimeRepo := sqlite.NewShowtimeRepo(db)
	holdStore := rediscache.NewSeatHoldStore(rdb)

	showtimeSvc := service.NewShowtimeService(showtimeRepo)
	bookingSvc := service.NewBookingService(showtimeRepo, holdStore, cfg.HoldTTL, cfg.MaxSeats)

	showtimeHandler := handler.NewShowtimeHandler(showtimeSvc)
	bookingHandler := handler.NewBookingHandler(bookingSvc)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler.NewRouter(showtimeHandler, bookingHandler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("listening on %s (hold TTL %s)", cfg.HTTPAddr, cfg.HoldTTL)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("server stopped")
}
