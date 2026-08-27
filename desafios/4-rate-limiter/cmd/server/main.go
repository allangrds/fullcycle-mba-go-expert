package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/configs"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/httpmiddleware"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/limiter"
	redisstorage "github.com/allangrds/fullcycle-mba-go-expert/desafios/4-rate-limiter/internal/storage/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := configs.Load(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer rdb.Close()

	// Strategy de persistência em uso: Redis. Para trocar por outro
	// mecanismo, basta implementar limiter.Storage em um novo pacote e
	// substituir esta instanciação — nada mais precisa mudar.
	storage := redisstorage.New(rdb)

	rateLimiter := limiter.New(storage, limiter.Config{
		IPMax:         cfg.RateLimitIPMax,
		TokenMax:      cfg.RateLimitTokenMax,
		Window:        cfg.Window(),
		BlockDuration: cfg.BlockDuration(),
	})

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(httpmiddleware.RateLimiter(rateLimiter))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	addr := ":" + cfg.WebServerPort
	log.Printf("rate limiter server running on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
