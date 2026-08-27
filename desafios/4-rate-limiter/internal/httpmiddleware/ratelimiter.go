// Package httpmiddleware adapta a regra de negócio de internal/limiter para
// o mundo HTTP. Não conhece Redis nem a lógica de contagem/bloqueio — só
// sabe montar uma limiter.Identity a partir do request e traduzir a
// resposta do limiter em um código HTTP.
package httpmiddleware

import (
	"net"
	"net/http"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/limiter"
)

const (
	apiKeyHeader   = "API_KEY"
	blockedMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

// RateLimiter constrói um middleware compatível com net/http (e com o
// router.Use do chi) que aplica o rate limiter em toda requisição.
func RateLimiter(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := extractIdentity(r)

			allowed, err := l.Allow(r.Context(), id)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(blockedMessage))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractIdentity(r *http.Request) limiter.Identity {
	if token := r.Header.Get(apiKeyHeader); token != "" {
		return limiter.Identity{Kind: limiter.KindToken, Value: token}
	}
	return limiter.Identity{Kind: limiter.KindIP, Value: clientIP(r)}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
