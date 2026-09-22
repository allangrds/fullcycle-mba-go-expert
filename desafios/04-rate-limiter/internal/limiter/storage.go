package limiter

import (
	"context"
	"time"
)

// Storage é o contrato do padrão Strategy para persistência do rate limiter.
// Qualquer mecanismo (Redis, memória, outro banco) pode ser usado desde que
// implemente estes três métodos.
type Storage interface {
	// Increment incrementa o contador da janela atual identificada por key e
	// retorna a contagem resultante. Se a chave ainda não existir, ela deve
	// ser criada com TTL igual a window.
	Increment(ctx context.Context, key string, window time.Duration) (int64, error)

	// IsBlocked informa se a chave de bloqueio está ativa no momento.
	IsBlocked(ctx context.Context, key string) (bool, error)

	// Block marca a chave como bloqueada pelo tempo informado. Deve ser
	// idempotente: se a chave já estiver bloqueada, não deve reiniciar o TTL.
	Block(ctx context.Context, key string, duration time.Duration) error
}
