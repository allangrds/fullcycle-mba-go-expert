package stress

import (
	"fmt"
	"sort"
	"time"
)

// Report consolida os resultados de uma execução de stress test.
type Report struct {
	TotalRequests int
	StatusCodes   map[int]int
	Errors        int
	Duration      time.Duration
}

// NewReport cria um Report vazio, pronto para receber resultados via Add.
func NewReport() Report {
	return Report{StatusCodes: make(map[int]int)}
}

// Add contabiliza um Result no relatório.
func (r *Report) Add(result Result) {
	r.TotalRequests++
	if result.Err != nil {
		r.Errors++
		return
	}
	r.StatusCodes[result.StatusCode]++
}

// Print imprime o relatório formatado no console.
func (r Report) Print() {
	fmt.Printf("Tempo total: %s\n", r.Duration)
	fmt.Printf("Total de requests: %d\n", r.TotalRequests)
	fmt.Printf("Requests com status 200: %d\n", r.StatusCodes[200])

	others := make([]int, 0, len(r.StatusCodes))
	for code := range r.StatusCodes {
		if code != 200 {
			others = append(others, code)
		}
	}

	if len(others) > 0 {
		sort.Ints(others)
		fmt.Println("Distribuição de outros status codes:")
		for _, code := range others {
			fmt.Printf("  %d: %d\n", code, r.StatusCodes[code])
		}
	}

	if r.Errors > 0 {
		fmt.Printf("Erros de rede/timeout: %d\n", r.Errors)
	}
}
