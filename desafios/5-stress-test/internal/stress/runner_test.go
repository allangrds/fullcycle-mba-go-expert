package stress_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/5-stress-test/internal/stress"
	"github.com/stretchr/testify/assert"
)

func TestRun_TotalRequestsExact(t *testing.T) {
	tests := []struct {
		name        string
		requests    int
		concurrency int
	}{
		{"concorrência divide o total", 20, 4},
		{"concorrência não divide o total", 7, 3},
		{"concorrência maior que o total", 3, 10},
		{"concorrência igual a um", 5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			report := stress.Run(context.Background(), server.URL, tt.requests, tt.concurrency)

			assert.Equal(t, tt.requests, report.TotalRequests)
			assert.Equal(t, tt.requests, report.StatusCodes[http.StatusOK])
		})
	}
}

func TestRun_StatusCodeDistribution(t *testing.T) {
	var counter int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt64(&counter, 1)
		switch {
		case n%5 == 0:
			w.WriteHeader(http.StatusInternalServerError)
		case n%3 == 0:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	const totalRequests = 30
	report := stress.Run(context.Background(), server.URL, totalRequests, 5)

	assert.Equal(t, totalRequests, report.TotalRequests)

	sum := report.Errors
	for _, count := range report.StatusCodes {
		sum += count
	}
	assert.Equal(t, totalRequests, sum)

	assert.Greater(t, report.StatusCodes[http.StatusOK], 0)
	assert.Greater(t, report.StatusCodes[http.StatusNotFound], 0)
	assert.Greater(t, report.StatusCodes[http.StatusInternalServerError], 0)
}

func TestRun_RespectsMaxConcurrency(t *testing.T) {
	const concurrency = 4

	var current, peak int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt64(&current, 1)
		for {
			m := atomic.LoadInt64(&peak)
			if n <= m || atomic.CompareAndSwapInt64(&peak, m, n) {
				break
			}
		}

		time.Sleep(10 * time.Millisecond)

		atomic.AddInt64(&current, -1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	report := stress.Run(context.Background(), server.URL, 40, concurrency)

	assert.Equal(t, 40, report.TotalRequests)
	assert.LessOrEqual(t, atomic.LoadInt64(&peak), int64(concurrency))
}

func TestRun_CountsNetworkErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	unreachableURL := server.URL
	server.Close() // fecha o servidor: as requisições vão falhar por conexão recusada

	report := stress.Run(context.Background(), unreachableURL, 5, 2)

	assert.Equal(t, 5, report.TotalRequests)
	assert.Equal(t, 5, report.Errors)
}
