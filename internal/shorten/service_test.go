package shorten_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"shortener/internal/model"
	"shortener/internal/shorten"
	"shortener/internal/storage/shortening"
	"testing"
)

func TestService_Shorten(t *testing.T) {
	t.Run("generates shortening for a given URL", func(t *testing.T) {
		var (
			svc   = shorten.NewService(shortening.NewInMemory())
			input = model.ShortenInput{RawURL: "https://google.com"}
		)
		short, err := svc.Shorten(context.Background(), input)
		require.NoError(t, err)

		require.NotEmpty(t, short.Identifier)
		assert.Equal(t, input.RawURL, short.OriginalUrl)
	})
	t.Run("uses custom identifier if provider", func(t *testing.T) {

	})
	t.Run("returns errors if identifier is already taken", func(t *testing.T) {

	})
}
