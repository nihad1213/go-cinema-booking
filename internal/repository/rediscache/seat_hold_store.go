package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"cinema-booking/internal/domain"
)

type SeatHoldStore struct{ rdb *redis.Client }

func NewSeatHoldStore(rdb *redis.Client) *SeatHoldStore {
	return &SeatHoldStore{rdb: rdb}
}

func key(showtimeID, seatID int64) string {
	return fmt.Sprintf("hold:{%d}:%d", showtimeID, seatID)
}

func keys(showtimeID int64, seatIDs []int64) []string {
	out := make([]string, len(seatIDs))
	for i, id := range seatIDs {
		out[i] = key(showtimeID, id)
	}
	return out
}

var holdScript = redis.NewScript(`
for i = 1, #KEYS do
  local owner = redis.call('GET', KEYS[i])
  if owner and owner ~= ARGV[1] then
    return 0
  end
end
for i = 1, #KEYS do
  redis.call('SET', KEYS[i], ARGV[1], 'PX', ARGV[2])
end
return 1
`)


var releaseScript = redis.NewScript(`
for i = 1, #KEYS do
  if redis.call('GET', KEYS[i]) == ARGV[1] then
    redis.call('DEL', KEYS[i])
  end
end
return 1
`)

func (s *SeatHoldStore) Hold(ctx context.Context, showtimeID int64, seatIDs []int64, userID string, ttl time.Duration) error {
	n, err := holdScript.Run(ctx, s.rdb, keys(showtimeID, seatIDs), userID, ttl.Milliseconds()).Int()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrSeatUnavailable
	}
	return nil
}

func (s *SeatHoldStore) Release(ctx context.Context, showtimeID int64, seatIDs []int64, userID string) error {
	return releaseScript.Run(ctx, s.rdb, keys(showtimeID, seatIDs), userID).Err()
}

func (s *SeatHoldStore) Holds(ctx context.Context, showtimeID int64, seatIDs []int64) (map[int64]string, error) {
	out := make(map[int64]string)
	if len(seatIDs) == 0 {
		return out, nil
	}

	vals, err := s.rdb.MGet(ctx, keys(showtimeID, seatIDs)...).Result()
	if err != nil {
		return nil, err
	}
	for i, v := range vals {
		if userID, ok := v.(string); ok {
			out[seatIDs[i]] = userID
		}
	}
	return out, nil
}
