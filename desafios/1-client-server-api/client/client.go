package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

var SERVER_TIMEOUT = 300 * time.Millisecond
var SERVER_URL = "http://localhost:8080/cotacao"
var FILE_TIMEOUT = 100 * time.Millisecond

func main() {
	// Criado context para requisição HTTP ao servidor
	ctx, cancel := context.WithTimeout(context.Background(), SERVER_TIMEOUT)
	defer cancel()

	// Criada request HTTP com contexto
	req, err := http.NewRequestWithContext(ctx, "GET", SERVER_URL, nil)
	if err != nil {
		log.Printf("Client | Erro ao criar request: %v", err)
		return
	}

	// Enviada request ao servidor local
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Client | Erro ao enviar request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Lida a resposta do servidor
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Client | Erro ao ler resposta do servidor: %v", err)
		return
	}

	// Parseada a resposta JSON
	var cotacao CotacaoResponse
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Printf("Client | Erro ao fazer parse da resposta: %v", err)
		return
	}

	// Criado contexto para salvar arquivo
	fileCtx, fileCancel := context.WithTimeout(context.Background(), FILE_TIMEOUT)
	defer fileCancel()

	// Salvada a cotação em arquivo
	err = salvarCotacaoArquivo(fileCtx, cotacao.Bid)
	if err != nil {
		log.Printf("Client | Erro ao salvar cotação em arquivo: %v", err)
		return
	}

	fmt.Printf("Cotação salva com sucesso: Dólar: %s\n", cotacao.Bid)
}

func salvarCotacaoArquivo(ctx context.Context, valor string) error {
	done := make(chan error, 1)

	go func() {
		conteudo := fmt.Sprintf("Dólar: %s", valor)
		err := os.WriteFile("cotacao.txt", []byte(conteudo), 0644)
		done <- err
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
