package defaults

import (
	"fmt"

	"github.com/rainproj/rain/pkg/context"
)

// Defaulter intergace
type Defaulter interface {
	fmt.Stringer

	Default(*context.Context) error
}

// Defaulters steps
var Defaulters = []Defaulter{}
