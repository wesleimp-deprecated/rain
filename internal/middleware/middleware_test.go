package middleware

import (
	"fmt"
	"testing"

	"github.com/rainproj/rain/pkg/context"
	"github.com/stretchr/testify/assert"
)

var ctx = &context.Context{}

func action(err error) Action {
	return func(ctx *context.Context) error {
		return err
	}
}

func TestLogging(t *testing.T) {
	assert.NoError(t, Log("foo", action(nil), DefaultInitialPadding)(ctx))
}

func TestError(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		assert.NoError(t, ErrHandler(action(nil))(ctx))
	})

	t.Run("with err", func(t *testing.T) {
		assert.Error(t, ErrHandler(action(fmt.Errorf("pipe errored")))(ctx))
	})
}
