package rest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLinks(t *testing.T) {
	linkURI := "http://example.com?page=2&limit=10"

	links := ParseLink(fmt.Sprintf(`<%s>; rel="next"`, linkURI))
	assert.Equal(t, 1, len(links))

	next := links.Next()
	assert.Equal(t, linkURI, next)
}
