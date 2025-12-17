package store

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/is_backend/services/admin/internal/config"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore() *RedisStore {
	return &RedisStore{
		ctx: context.Background(),
	}
}

func (s *RedisStore) Init(cfg *config.Config) {
	s.client = redis.NewClient(&redis.Options{
		// Addr:     cfg.Redis.Addr,
		// Password: cfg.Redis.Pass,
		// DB:       cfg.Redis.DB,
	})
}

func (s *RedisStore) Connect() error {
	_, err := s.client.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStore) GetKey(key string) (valuer string, err error) {
	val, err := s.client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (s *RedisStore) SetKey(key string, value string, ttl time.Duration) error {
	return s.client.Set(key, value, ttl).Err()
}
