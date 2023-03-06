package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

func NewStorage(ctx context.Context, address string) (Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if cmd := client.Ping(ctx); cmd.Err() != nil {
		return Storage{}, cmd.Err()
	}

	return Storage{
		client: client,
	}, nil
}

func (s *Storage) LockResource(ctx context.Context, resource, randomValue string, ttl time.Duration) (bool, error) {
	cmd := s.client.SetNX(ctx, resource, randomValue, ttl)
	return cmd.Result()
}

func (s *Storage) UnlockResource(ctx context.Context, resource, randomValue string) error {
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)

	keys := []string{resource}
	values := []interface{}{randomValue}
	cmd := script.Run(ctx, s.client, keys, values)
	return cmd.Err()
}
