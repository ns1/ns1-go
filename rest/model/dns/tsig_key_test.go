package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTsigKey(t *testing.T) {
	tsigKey := NewTsigKey("myTsigKey", "myAlgorithm", "mySecret")

	assert.Equal(t, "myTsigKey", tsigKey.Name)
	assert.Equal(t, "myAlgorithm", tsigKey.Algorithm)
	assert.Equal(t, "mySecret", tsigKey.Secret)
}
