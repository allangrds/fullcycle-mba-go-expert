package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/5-stress-test/internal/stress"
	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
)

// rootCmd representa o comando base do CLI de stress test.
var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "CLI de stress test HTTP",
	Long: `Dispara um número configurável de requisições HTTP concorrentes
contra uma URL e gera, ao final, um relatório com o tempo total gasto,
o total de requests realizados e a distribuição dos status codes.`,
	RunE: run,
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 0, "número total de requisições (obrigatório)")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 0, "número de chamadas simultâneas (obrigatório)")

	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.MarkFlagRequired("concurrency")
}

// Execute roda o comando raiz do CLI. É chamado por main.main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if requests <= 0 {
		return fmt.Errorf("--requests deve ser maior que zero")
	}
	if concurrency <= 0 {
		return fmt.Errorf("--concurrency deve ser maior que zero")
	}

	report := stress.Run(context.Background(), url, requests, concurrency)
	report.Print()

	return nil
}
