package user_entity

import (
	"context"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/8-fechamento-automatico-leiloes/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(
		ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
