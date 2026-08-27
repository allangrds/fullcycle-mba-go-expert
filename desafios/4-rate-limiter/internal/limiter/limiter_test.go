package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/limiter"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func newLimiter(ipMax, tokenMax int64, window, block time.Duration) *limiter.Limiter {
	return limiter.New(memory.New(), limiter.Config{
		IPMax:         ipMax,
		TokenMax:      tokenMax,
		Window:        window,
		BlockDuration: block,
	})
}

func TestLimiter_Allow_WithinAndOverLimit(t *testing.T) {
	tests := []struct {
		name     string
		kind     limiter.Kind
		max      int64
		requests int
		want     []bool
	}{
		{
			name:     "ip within limit",
			kind:     limiter.KindIP,
			max:      3,
			requests: 3,
			want:     []bool{true, true, true},
		},
		{
			name:     "ip exceeds limit on the 4th request",
			kind:     limiter.KindIP,
			max:      3,
			requests: 4,
			want:     []bool{true, true, true, false},
		},
		{
			name:     "token exceeds limit on the 6th request",
			kind:     limiter.KindToken,
			max:      5,
			requests: 6,
			want:     []bool{true, true, true, true, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l *limiter.Limiter
			if tt.kind == limiter.KindIP {
				l = newLimiter(tt.max, 1000, time.Minute, time.Minute)
			} else {
				l = newLimiter(1000, tt.max, time.Minute, time.Minute)
			}

			id := limiter.Identity{Kind: tt.kind, Value: "identity-" + tt.name}

			got := make([]bool, 0, tt.requests)
			for i := 0; i < tt.requests; i++ {
				allowed, err := l.Allow(context.Background(), id)
				assert.NoError(t, err)
				got = append(got, allowed)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLimiter_StaysBlockedDuringBlockWindow(t *testing.T) {
	// Window bem menor que BlockDuration, como no mundo real (ex.: janela de
	// 1s e bloqueio de 5 minutos): assim, quando o bloqueio expira, a chave
	// de contagem também já expirou e a contagem recomeça do zero.
	l := newLimiter(1, 1000, 100*time.Millisecond, 200*time.Millisecond)
	id := limiter.Identity{Kind: limiter.KindIP, Value: "1.2.3.4"}

	allowed, err := l.Allow(context.Background(), id)
	assert.NoError(t, err)
	assert.True(t, allowed)

	allowed, err = l.Allow(context.Background(), id)
	assert.NoError(t, err)
	assert.False(t, allowed, "second request should exceed the limit and trigger a block")

	allowed, err = l.Allow(context.Background(), id)
	assert.NoError(t, err)
	assert.False(t, allowed, "still inside the block window")

	time.Sleep(300 * time.Millisecond)

	allowed, err = l.Allow(context.Background(), id)
	assert.NoError(t, err)
	assert.True(t, allowed, "block window has expired")
}

// TestLimiter_TokenPrecedesIP prova a "regra de ouro" do desafio: o mesmo
// valor usado como IP e como Token tem limites e bloqueios totalmente
// independentes -- exceder o limite do IP não afeta o Token.
func TestLimiter_TokenPrecedesIP(t *testing.T) {
	l := newLimiter(2, 10, time.Minute, time.Minute)
	same := "192.168.0.10"

	ipID := limiter.Identity{Kind: limiter.KindIP, Value: same}
	tokenID := limiter.Identity{Kind: limiter.KindToken, Value: same}

	for i := 0; i < 2; i++ {
		allowed, err := l.Allow(context.Background(), ipID)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, err := l.Allow(context.Background(), ipID)
	assert.NoError(t, err)
	assert.False(t, allowed, "IP should be blocked after exceeding RATE_LIMIT_IP_MAX")

	for i := 0; i < 10; i++ {
		allowed, err := l.Allow(context.Background(), tokenID)
		assert.NoError(t, err)
		assert.True(t, allowed, "token limit should not be affected by the blocked IP")
	}
}

// TestLimiter_ConcurrentRequests garante que, mesmo com várias goroutines
// disparando ao mesmo tempo, exatamente IPMax requisições são aceitas.
func TestLimiter_ConcurrentRequests(t *testing.T) {
	const ipMax = 20
	const totalRequests = 100

	l := newLimiter(ipMax, 1000, time.Minute, time.Minute)
	id := limiter.Identity{Kind: limiter.KindIP, Value: "10.0.0.1"}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		allowed int
	)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, err := l.Allow(context.Background(), id)
			assert.NoError(t, err)
			if ok {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, ipMax, allowed)
}
