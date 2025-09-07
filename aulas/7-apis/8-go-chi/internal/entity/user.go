package entity

import (
	"github.com/allangrds/fullcycle-mba-go-expert/aulas/7-apis/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

// Usando o - para omitir o campo da serialização JSON
type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

/*
Pq return &User ao invés de User?
Para evitar a cópia desnecessária de dados e garantir que as alterações na estrutura sejam refletidas em todas as referências
*/
func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, err
}

func (u *User) IsPasswordValid(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
