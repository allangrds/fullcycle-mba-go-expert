package stress

import (
	"context"
	"net/http"
	"sync"
	"time"
)

const requestTimeout = 10 * time.Second

// Result é o resultado de uma única requisição HTTP.
type Result struct {
	StatusCode int
	Err        error
}

// Run dispara exatamente `requests` requisições HTTP GET para `url`,
// distribuídas entre `concurrency` workers, e devolve o relatório
// consolidado ao final.
func Run(ctx context.Context, url string, requests, concurrency int) Report {
	start := time.Now()

	jobs := make(chan struct{}, requests)
	for i := 0; i < requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	results := make(chan Result, requests)
	client := &http.Client{}

	var wg sync.WaitGroup
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				results <- doRequest(ctx, client, url)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	report := NewReport()
	for result := range results {
		report.Add(result)
	}
	report.Duration = time.Since(start)

	return report
}

func doRequest(ctx context.Context, client *http.Client, url string) Result {
	reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return Result{Err: err}
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{Err: err}
	}
	defer resp.Body.Close()

	return Result{StatusCode: resp.StatusCode}
}
