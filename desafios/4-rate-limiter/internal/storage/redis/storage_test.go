package redis_test

import (
	"context"
	"os"
	"testing"
	"time"

	redisstorage "github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/storage/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestClient conecta a um Redis real (REDIS_HOST/REDIS_PORT, com padrão
// localhost:6379). Se o Redis não estiver acessível, o teste é pulado em vez
// de falhar -- ele é pensado para rodar via "docker compose run --rm test",
// onde o serviço redis do compose está sempre disponível.
func newTestClient(t *testing.T) *goredis.Client {
	t.Helper()

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:        host + ":" + port,
		MaxRetries:  -1,
		DialTimeout: 300 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("redis not reachable at %s:%s, skipping integration test: %v", host, port, err)
	}

	return client
}

func TestStorage_IncrementAndBlock(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	s := redisstorage.New(client)
	ctx := context.Background()

	countKey := "test:rate_limiter:count:" + t.Name()
	blockKey := "test:rate_limiter:block:" + t.Name()
	defer client.Del(ctx, countKey, blockKey)

	count, err := s.Increment(ctx, countKey, time.Second)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	count, err = s.Increment(ctx, countKey, time.Second)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	blocked, err := s.IsBlocked(ctx, blockKey)
	require.NoError(t, err)
	assert.False(t, blocked)

	require.NoError(t, s.Block(ctx, blockKey, 200*time.Millisecond))

	blocked, err = s.IsBlocked(ctx, blockKey)
	require.NoError(t, err)
	assert.True(t, blocked)

	time.Sleep(250 * time.Millisecond)

	blocked, err = s.IsBlocked(ctx, blockKey)
	require.NoError(t, err)
	assert.False(t, blocked, "block TTL should have expired in Redis")
}

func TestStorage_IncrementResetsAfterWindowExpires(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	s := redisstorage.New(client)
	ctx := context.Background()

	key := "test:rate_limiter:count:" + t.Name()
	defer client.Del(ctx, key)

	_, err := s.Increment(ctx, key, 200*time.Millisecond)
	require.NoError(t, err)

	time.Sleep(250 * time.Millisecond)

	count, err := s.Increment(ctx, key, 200*time.Millisecond)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Redis TTL should have expired the key")
}
