package logone

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	t.Run("scenario", func(t *testing.T) {
		l := NewLoggerDefault()
		ctx := context.Background()
		ctx = NewContextWithLogger(ctx, l)
		nl, ok := LoggerFromContext(ctx)
		assert.Exactly(t, true, ok)
		assert.Exactly(t, l, nl)
	})
}
