package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	assert.NotEmpty(t, Step{}.String())
}
