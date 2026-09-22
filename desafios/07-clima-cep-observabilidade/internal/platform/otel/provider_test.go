package otelprovider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProvider(t *testing.T) {
	shutdown, err := NewProvider(context.Background(), "test-service", "localhost:4317")

	assert.NoError(t, err)
	assert.NotNil(t, shutdown)
	assert.NoError(t, shutdown(context.Background()))
}
