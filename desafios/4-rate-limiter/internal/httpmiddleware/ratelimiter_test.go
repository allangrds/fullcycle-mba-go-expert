package httpmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/httpmiddleware"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/limiter"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

const blockedBody = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func newHandler(ipMax, tokenMax int64) http.Handler {
	l := limiter.New(memory.New(), limiter.Config{
		IPMax:         ipMax,
		TokenMax:      tokenMax,
		Window:        time.Minute,
		BlockDuration: time.Minute,
	})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("next-called"))
	})

	return httpmiddleware.RateLimiter(l)(next)
}

func doRequest(handler http.Handler, remoteAddr, apiKey string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = remoteAddr
	if apiKey != "" {
		req.Header.Set("API_KEY", apiKey)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestRateLimiter_AllowsWithinLimit(t *testing.T) {
	handler := newHandler(2, 10)

	for i := 0; i < 2; i++ {
		rec := doRequest(handler, "1.2.3.4:5555", "")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "next-called", rec.Body.String())
	}
}

func TestRateLimiter_BlocksWithExactBodyAndStatus(t *testing.T) {
	handler := newHandler(1, 10)

	rec := doRequest(handler, "1.2.3.4:5555", "")
	assert.Equal(t, http.StatusOK, rec.Code)

	rec = doRequest(handler, "1.2.3.4:5555", "")
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Equal(t, blockedBody, rec.Body.String())
}

// TestRateLimiter_TokenPrecedesIP prova a regra de ouro na camada HTTP: uma
// requisição do mesmo IP já bloqueado, mas com um API_KEY válido, continua
// sendo aceita porque usa o limite (maior) do Token.
func TestRateLimiter_TokenPrecedesIP(t *testing.T) {
	handler := newHandler(1, 10)

	rec := doRequest(handler, "9.9.9.9:1111", "")
	assert.Equal(t, http.StatusOK, rec.Code)

	rec = doRequest(handler, "9.9.9.9:1111", "")
	assert.Equal(t, http.StatusTooManyRequests, rec.Code, "IP should now be blocked")

	rec = doRequest(handler, "9.9.9.9:1111", "any-token")
	assert.Equal(t, http.StatusOK, rec.Code, "token limit takes precedence over the blocked IP")
}
