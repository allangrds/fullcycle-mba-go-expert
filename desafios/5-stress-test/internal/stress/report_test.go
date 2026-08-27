package stress_test

import (
	"testing"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/5-stress-test/internal/stress"
	"github.com/stretchr/testify/assert"
)

func TestReport_Add(t *testing.T) {
	tests := []struct {
		name            string
		results         []stress.Result
		wantTotal       int
		wantStatusCodes map[int]int
		wantErrors      int
	}{
		{
			name: "só sucessos",
			results: []stress.Result{
				{StatusCode: 200},
				{StatusCode: 200},
				{StatusCode: 200},
			},
			wantTotal:       3,
			wantStatusCodes: map[int]int{200: 3},
			wantErrors:      0,
		},
		{
			name: "mistura de status codes",
			results: []stress.Result{
				{StatusCode: 200},
				{StatusCode: 404},
				{StatusCode: 500},
				{StatusCode: 200},
			},
			wantTotal:       4,
			wantStatusCodes: map[int]int{200: 2, 404: 1, 500: 1},
			wantErrors:      0,
		},
		{
			name: "com erros de rede",
			results: []stress.Result{
				{StatusCode: 200},
				{Err: assert.AnError},
				{Err: assert.AnError},
			},
			wantTotal:       3,
			wantStatusCodes: map[int]int{200: 1},
			wantErrors:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := stress.NewReport()
			for _, result := range tt.results {
				report.Add(result)
			}

			assert.Equal(t, tt.wantTotal, report.TotalRequests)
			assert.Equal(t, tt.wantErrors, report.Errors)
			for code, count := range tt.wantStatusCodes {
				assert.Equal(t, count, report.StatusCodes[code])
			}
		})
	}
}

func TestReport_Print_DoesNotPanic(t *testing.T) {
	report := stress.NewReport()
	report.Add(stress.Result{StatusCode: 200})
	report.Add(stress.Result{StatusCode: 404})
	report.Add(stress.Result{Err: assert.AnError})
	report.Duration = 1500 * time.Millisecond

	assert.NotPanics(t, func() {
		report.Print()
	})
}
