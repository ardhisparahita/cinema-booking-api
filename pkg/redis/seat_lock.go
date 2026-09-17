package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type SeatLocker interface {
	LockSeats(ctx context.Context, showtimeID uint, seatIDs []uint, value string, ttl time.Duration) (bool, error)
	UnlockSeats(ctx context.Context, showtimeID uint, seatIDs []uint, value string) error
}

type SeatLockerImpl struct {
	Client *goredis.Client
}

func NewSeatLocker(client *goredis.Client) SeatLocker {
	return &SeatLockerImpl{
		Client: client,
	}
}

func seatLockKey(showtimeID, seatID uint) string {
	return fmt.Sprintf(
		"booking:showtime:%d:seat:%d",
		showtimeID,
		seatID,
	)
}

var lockSeatsScript = goredis.NewScript(`
for _, key in ipairs(KEYS) do 
	if redis.call("EXIST", key) == 1 then
		return 0
	end
end

for _, key in ipairs(KEYS) do
	redis.call("SET", key, ARGV[1], "PX", ARGV[2])
end

return 1
`)

var unlockSeatsScript = goredis.NewScript(`
local deleted = 0

for _, key in ipairs(KEYS) do
	if redis.call("GET", key) == ARGV[1] then
		redis.call("DEL", key)
		deleted = deleted + 1
	end
end

return deleted
`)

func (s *SeatLockerImpl) LockSeats(ctx context.Context, showtimeID uint, seatIDs []uint, value string, ttl time.Duration) (bool, error) {
	if len(seatIDs) == 0 {
		return false, nil
	}

	if value == "" {
		return false, fmt.Errorf("lock value is required")
	}
	if ttl <= 0 {
		return false, fmt.Errorf("lock ttl must be greater than zero")
	}

	keys := make([]string, 0, len(seatIDs))

	for _, seatID := range seatIDs {
		keys = append(keys, seatLockKey(showtimeID, seatID))
	}

	ttlMillis := ttl.Milliseconds()

	result, err := lockSeatsScript.Run(
		ctx,
		s.Client,
		keys,
		value,
		strconv.FormatInt(ttlMillis, 10),
	).Int()

	if err != nil {
		return false, fmt.Errorf("failed to lock seat: %w", err)
	}

	return result == 1, nil
}

func (s *SeatLockerImpl) UnlockSeats(ctx context.Context, showtimeID uint, seatIDs []uint, value string) error {
	if len(seatIDs) == 0 {
		return nil
	}

	if value == "" {
		return fmt.Errorf("lock value is required")
	}

	keys := make([]string, 0, len(seatIDs))

	for _, seatID := range seatIDs {
		keys = append(keys, seatLockKey(showtimeID, seatID))
	}

	if _, err := unlockSeatsScript.Run(
		ctx,
		s.Client,
		keys,
		value,
	).Int(); err != nil {
		return fmt.Errorf("failed to unlock seats: %w", err)
	}

	return nil
}
