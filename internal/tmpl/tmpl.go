package tmpl

import (
	"bytes"
	"text/template"
	"time"

	"github.com/rainproj/rain/pkg/context"
)

// Template holds data that can be applied to a template string.
type Template struct {
	fields Fields
}

// Fields that will be available to the template engine.
type Fields map[string]interface{}

const (
	name      = "Name"
	env       = "Env"
	date      = "Date"
	timestamp = "Timestamp"
	version   = "Version"
)

// New Template.
func New(ctx *context.Context) *Template {
	return &Template{
		fields: Fields{
			name:      ctx.Config.ProjectName,
			env:       ctx.Env,
			version:   ctx.Config.Version,
			date:      ctx.Date.UTC().Format(time.RFC3339),
			timestamp: ctx.Date.UTC().Unix(),
		},
	}
}

// WithEnv overrides template's env field with the given environment map.
func (t *Template) WithEnv(e map[string]string) *Template {
	t.fields[env] = e
	return t
}

// WithExtraFields allows to add new more custom fields to the template.
// It will override fields with the same name.
func (t *Template) WithExtraFields(f Fields) *Template {
	for k, v := range f {
		t.fields[k] = v
	}
	return t
}

// Apply applies the given string against the Fields stored in the template.
func (t *Template) Apply(s string) (string, error) {
	var out bytes.Buffer
	tmpl, err := template.New("tmpl").
		Option("missingkey=error").
		Parse(s)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&out, t.fields)
	return out.String(), err
}
