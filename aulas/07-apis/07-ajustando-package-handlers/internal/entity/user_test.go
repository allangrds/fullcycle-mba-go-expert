package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("Fulano de Tal", "joao@detal.com.br", "123456")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.NotEqual(t, "123456", user.Password)
	assert.Equal(t, "Fulano de Tal", user.Name)
	assert.Equal(t, "joao@detal.com.br", user.Email)
}

func TestIsPasswordValid(t *testing.T) {
	user, err := NewUser("Fulano de Tal", "joao@detal.com.br", "123456")

	assert.Nil(t, err)
	assert.True(t, user.IsPasswordValid("123456"))
	assert.False(t, user.IsPasswordValid("senhaerrada"))
}
