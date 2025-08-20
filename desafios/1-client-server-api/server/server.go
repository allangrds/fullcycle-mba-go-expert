package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CotacaoResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type CotacaoResult struct {
	Bid string `json:"bid"`
}

type Cotacao struct {
	ID        uint      `gorm:"primaryKey"`
	Bid       string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

var API_TIMEOUT = 200 * time.Millisecond
var API_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
var DB_TIMEOUT = 10 * time.Millisecond

func main() {
	http.HandleFunc("/cotacao", handler)
	log.Println("Server | Rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Iniciado banco de dados com GORM
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		log.Printf("Server | Erro ao conectar com o banco de dados: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Feito auto-migrate para criar a tabela
	db.AutoMigrate(&Cotacao{})

	// Criado context para chamar a API
	ctx, cancel := context.WithTimeout(r.Context(), API_TIMEOUT)
	defer cancel()

	// Criada a request HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", API_URL, nil)
	if err != nil {
		log.Printf("Server | Erro ao criar request: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Feita a chamada da API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Server | Erro ao chamar a API: %v", err)
		http.Error(w, "Erro ao buscar cotação", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lida a resposta da API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Server | Erro ao ler resposta da API: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Feito o parse da resposta JSON
	var cotacaoResp CotacaoResponse
	err = json.Unmarshal(body, &cotacaoResp)
	if err != nil {
		log.Printf("Server | Erro ao fazer parse da resposta: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	// Criado context para salvar no banco
	dbCtx, dbCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer dbCancel()

	// Salva no banco usando GORM
	cotacao := Cotacao{
		Bid: cotacaoResp.USDBRL.Bid,
	}
	result := db.WithContext(dbCtx).Create(&cotacao)
	if result.Error != nil {
		log.Printf("Server | Erro ao salvar no banco: %v", result.Error)
		http.Error(w, "Erro ao salvar cotação", http.StatusInternalServerError)
		return
	}

	// Monta resposta para o client
	response := CotacaoResult{
		Bid: cotacaoResp.USDBRL.Bid,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
