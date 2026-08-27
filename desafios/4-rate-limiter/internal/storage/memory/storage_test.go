package memory_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestStorage_Increment(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	count, err := s.Increment(ctx, "key", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	count, err = s.Increment(ctx, "key", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestStorage_IncrementResetsAfterWindowExpires(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	_, err := s.Increment(ctx, "key", 50*time.Millisecond)
	assert.NoError(t, err)

	time.Sleep(60 * time.Millisecond)

	count, err := s.Increment(ctx, "key", 50*time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count, "counter should restart after the window expires")
}

func TestStorage_BlockIsIdempotent(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	assert.NoError(t, s.Block(ctx, "blocked", 50*time.Millisecond))

	blocked, err := s.IsBlocked(ctx, "blocked")
	assert.NoError(t, err)
	assert.True(t, blocked)

	// Chamar Block de novo não deve reiniciar o TTL do bloqueio original.
	time.Sleep(30 * time.Millisecond)
	assert.NoError(t, s.Block(ctx, "blocked", 50*time.Millisecond))

	time.Sleep(30 * time.Millisecond) // 60ms desde o primeiro Block

	blocked, err = s.IsBlocked(ctx, "blocked")
	assert.NoError(t, err)
	assert.False(t, blocked, "the original 50ms TTL should have expired")
}

func TestStorage_ConcurrentIncrement(t *testing.T) {
	s := memory.New()
	ctx := context.Background()

	const goroutines = 50
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.Increment(ctx, "shared", time.Minute)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	count, err := s.Increment(ctx, "shared", time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, int64(goroutines+1), count)
}
