// Package memory implementa a Strategy de persistência do rate limiter em
// memória. Serve para testes rápidos (sem depender de um Redis real) e como
// prova de que a persistência pode ser trocada sem alterar internal/limiter.
package memory

import (
	"context"
	"sync"
	"time"
)

type counter struct {
	value   int64
	expires time.Time
}

// Storage implementa limiter.Storage guardando contadores e bloqueios em
// mapas protegidos por um sync.Mutex.
type Storage struct {
	mu      sync.Mutex
	counts  map[string]counter
	blocked map[string]time.Time
}

func New() *Storage {
	return &Storage{
		counts:  make(map[string]counter),
		blocked: make(map[string]time.Time),
	}
}

func (s *Storage) Increment(ctx context.Context, key string, window time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	c, ok := s.counts[key]
	if !ok || now.After(c.expires) {
		c = counter{value: 0, expires: now.Add(window)}
	}
	c.value++
	s.counts[key] = c

	return c.value, nil
}

func (s *Storage) IsBlocked(ctx context.Context, key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	until, ok := s.blocked[key]
	if !ok {
		return false, nil
	}
	if time.Now().After(until) {
		delete(s.blocked, key)
		return false, nil
	}
	return true, nil
}

func (s *Storage) Block(ctx context.Context, key string, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.blocked[key]; ok {
		return nil
	}
	s.blocked[key] = time.Now().Add(duration)
	return nil
}
