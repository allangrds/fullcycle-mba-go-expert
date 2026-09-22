package limiter

import (
	"context"
	"time"
)

// Kind identifica se uma Identity representa um IP ou um Token de acesso.
type Kind string

const (
	KindIP    Kind = "ip"
	KindToken Kind = "token"
)

// Identity representa quem está fazendo a requisição: um IP ou um Token.
type Identity struct {
	Kind  Kind
	Value string
}

// Config contém os parâmetros de negócio do rate limiter, todos vindos de
// variáveis de ambiente (ver configs.Conf).
type Config struct {
	IPMax         int64
	TokenMax      int64
	Window        time.Duration
	BlockDuration time.Duration
}

// Limiter concentra a regra de negócio do rate limiter. Não conhece HTTP nem
// a tecnologia de persistência usada — apenas a interface Storage.
type Limiter struct {
	storage Storage
	cfg     Config
}

func New(storage Storage, cfg Config) *Limiter {
	return &Limiter{storage: storage, cfg: cfg}
}

// Allow decide se uma requisição da identidade informada pode prosseguir.
// Token e IP usam chaves e limites totalmente independentes, o que garante
// a precedência: o limite de Token nunca é afetado pelo estado do IP.
func (l *Limiter) Allow(ctx context.Context, id Identity) (bool, error) {
	blockKey := "rate_limiter:block:" + string(id.Kind) + ":" + id.Value

	blocked, err := l.storage.IsBlocked(ctx, blockKey)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	max := l.cfg.IPMax
	if id.Kind == KindToken {
		max = l.cfg.TokenMax
	}

	countKey := "rate_limiter:count:" + string(id.Kind) + ":" + id.Value
	count, err := l.storage.Increment(ctx, countKey, l.cfg.Window)
	if err != nil {
		return false, err
	}

	if count > max {
		if err := l.storage.Block(ctx, blockKey, l.cfg.BlockDuration); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
