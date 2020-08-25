package tmpl

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rainproj/rain/pkg/config"
	"github.com/rainproj/rain/pkg/context"
)

func TestEnv(t *testing.T) {
	testCases := []struct {
		desc string
		in   string
		out  string
	}{
		{
			desc: "with env",
			in:   "{{ .Env.FOO }}",
			out:  "BAR",
		},
		{
			desc: "with env",
			in:   "{{ .Env.BAR }}",
			out:  "",
		},
	}
	var ctx = context.New(config.Config{})
	ctx.Env = map[string]string{
		"FOO": "BAR",
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			out, _ := New(ctx).Apply(tC.in)
			assert.Equal(t, tC.out, out)
		})
	}
}

func TestInvalidTemplate(t *testing.T) {
	ctx := context.New(config.Config{})
	_, err := New(ctx).Apply("{{{.Foo}")
	assert.EqualError(t, err, "template: tmpl:1: unexpected \"{\" in command")
}

func TestEnvNotFound(t *testing.T) {
	var ctx = context.New(config.Config{})
	result, err := New(ctx).Apply("{{.Env.FOO}}")
	assert.Empty(t, result)
	assert.EqualError(t, err, `template: tmpl:1:6: executing "tmpl" at <.Env.FOO>: map has no entry for key "FOO"`)
}

func TestWithExtraFields(t *testing.T) {
	var ctx = context.New(config.Config{})
	out, _ := New(ctx).WithExtraFields(Fields{
		"MyCustomField": "foo",
	}).Apply("{{ .MyCustomField }}")
	assert.Equal(t, "foo", out)
}
