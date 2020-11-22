package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sergeyzalunin/go-shortener/internal/shortener"
)

// redisRepository is an another implementation of repository pattern.
type redisRepository struct {
	redisTimeout time.Duration
	client       *redis.Client
}

func newRedisClient(redisURL string, redisTimeout int) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.newRedisClient.ParseURL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(redisTimeout)*time.Second)
	defer cancel()

	client := redis.NewClient(opts)
	_, err = client.Ping(ctx).Result()

	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.newRedisClient.Ping")
	}

	return client, nil
}

// NewRedisRepository is a constructor for an implementation of RedirectRepository.
func NewRedisRepository(redisURL string, redisTimeout int) (shortener.RedirectRepository, error) {
	client, err := newRedisClient(redisURL, redisTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.NewRedisRepository")
	}

	repo := &redisRepository{
		redisTimeout: time.Duration(redisTimeout) * time.Second,
		client:       client,
	}

	return repo, nil
}

// generateKey creates the key for a redis.
// The key looks like "redirect:code",
// where "code" is shortened code of redirection.
func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

// Find looks for shortened url in the redis by provided code.
func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.redisTimeout)
	defer cancel()

	key := r.generateKey(code)

	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.Find.HGetAll")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, fmt.Sprintf("repository.redis.Find %v", len(data)))
	}

	createdAt, err := strconv.ParseInt(data[shortener.CreatedAtField], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.Find.ParseInt")
	}

	redirect := &shortener.Redirect{
		Code:      data[shortener.CodeField],
		URL:       data[shortener.URLField],
		CreatedAt: createdAt,
	}

	return redirect, nil
}

// Store saves the shortened url in Redis.
func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		shortener.CodeField:      redirect.Code,
		shortener.URLField:       redirect.URL,
		shortener.CreatedAtField: redirect.CreatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.redisTimeout)
	defer cancel()

	_, err := r.client.HSet(ctx, key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.redis.Store.HMSet")
	}

	return nil
}
