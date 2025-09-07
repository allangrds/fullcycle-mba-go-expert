package entity

import (
	"errors"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/pkg/entity"
)

var (
	ErrIDIsRequired    = errors.New("id is required")
	ErrInvalidID       = errors.New("invalid id")
	ErrNameIsRequired  = errors.New("name is required")
	ErrPriceIsRequired = errors.New("price is required")
	ErrInvalidPrice    = errors.New("invalid price")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Product) Validate() error {
	// Verificar se o ID está vazio
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}

	// Verificar se o ID é válido
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidID
	}

	// Verificar se o nome está vazio
	if p.Name == "" {
		return ErrNameIsRequired
	}

	// Verificar se o preço é zero
	if p.Price == 0 {
		return ErrPriceIsRequired
	}

	// Verificar se o preço é menor ou igual a zero
	if p.Price < 0 {
		return ErrInvalidPrice
	}

	return nil
}

func NewProduct(name string, price int) (*Product, error) {
	product := &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := product.Validate()
	if err != nil {
		return nil, err
	}

	return product, nil
}
