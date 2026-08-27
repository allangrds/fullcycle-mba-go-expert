// Package redis implementa a Strategy de persistência do rate limiter usando
// Redis. É a implementação usada em produção pelo desafio.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Storage implementa limiter.Storage sobre um cliente Redis. Usa apenas
// comandos atômicos nativos do Redis (SETNX + INCR), sem scripts Lua:
//   - SETNX cria a chave com o valor inicial e já define o TTL em uma única
//     operação atômica: se ela não existir, este é o primeiro hit da janela.
//   - Se SETNX falhar (a chave já existe), INCR soma 1 de forma atômica.
type Storage struct {
	client *redis.Client
}

func New(client *redis.Client) *Storage {
	return &Storage{client: client}
}

func (s *Storage) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	created, err := s.client.SetNX(ctx, key, 1, window).Result()
	if err != nil {
		return 0, err
	}
	if created {
		return 1, nil
	}
	return s.client.Incr(ctx, key).Result()
}

func (s *Storage) IsBlocked(ctx context.Context, key string) (bool, error) {
	n, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (s *Storage) Block(ctx context.Context, key string, duration time.Duration) error {
	return s.client.SetNX(ctx, key, 1, duration).Err()
}
