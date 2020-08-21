package context

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wesleimp/rain/pkg/config"
)

func TestNew(t *testing.T) {
	assert.NoError(t, os.Setenv("FOO", "NOT BAR"))
	assert.NoError(t, os.Setenv("BAR", "1"))
	var ctx = New(config.Config{
		Env: []string{
			"FOO=BAR",
		},
	})
	assert.Equal(t, "BAR", ctx.Env["FOO"])
	assert.Equal(t, "1", ctx.Env["BAR"])
	assert.Equal(t, 4, ctx.Parallelism)
}

func TestNewWithTimeout(t *testing.T) {
	ctx, cancel := NewWithTimeout(config.Config{}, time.Second)
	assert.NotEmpty(t, ctx.Env)
	assert.Equal(t, 4, ctx.Parallelism)
	cancel()
	<-ctx.Done()
	assert.EqualError(t, ctx.Err(), `context canceled`)
}
