package context

import (
	ctx "context"
	"os"
	"strings"
	"time"

	"github.com/wesleimp/rain/pkg/config"
)

// Context carriers all data through the pipes
type Context struct {
	ctx.Context
	Config      config.Config
	Env         Env
	Date        time.Time
	Timeout     time.Duration
	Parallelism int
	RmDist      bool
}

// Env is the environment variables.
type Env map[string]string

// Strings returns the current environment as a list of strings
func (e Env) Strings() []string {
	var result = make([]string, 0, len(e))
	for k, v := range e {
		result = append(result, k+"="+v)
	}
	return result
}

// New context.
func New(config config.Config) *Context {
	return Wrap(ctx.Background(), config)
}

// NewWithTimeout new context with the given timeout.
func NewWithTimeout(config config.Config, timeout time.Duration) (*Context, ctx.CancelFunc) {
	ctx, cancel := ctx.WithTimeout(ctx.Background(), timeout)
	return Wrap(ctx, config), cancel
}

// Wrap wraps an existing context.
func Wrap(ctx ctx.Context, config config.Config) *Context {
	return &Context{
		Context:     ctx,
		Config:      config,
		Env:         splitEnv(append(os.Environ(), config.Env...)),
		Parallelism: 4,
		Date:        time.Now(),
	}
}

func splitEnv(env []string) map[string]string {
	r := map[string]string{}
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		k := p[0]
		v := p[1]
		r[k] = v
	}

	return r
}
