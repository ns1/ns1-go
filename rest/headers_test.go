package rest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders_ParseLinks(t *testing.T) {
	linkURI := "http://example.com?page=2&limit=10"

	links := ParseLink(fmt.Sprintf(`<%s>; rel="next"`, linkURI), false)
	assert.Equal(t, 1, len(links))

	next := links.Next()
	assert.Equal(t, linkURI, next)
}

func TestHeaders_ParseLinksWithForceHTTPS(t *testing.T) {
	// It should replace HTTP with HTTPS in the urls
	linkURI := "http://example.com?page=2&limit=10"

	links := ParseLink(fmt.Sprintf(`<%s>; rel="next"`, linkURI), true)
	assert.Equal(t, 1, len(links))

	next := links.Next()
	assert.Equal(t, strings.Replace(linkURI, "http://", "https://", 1), next)
}
