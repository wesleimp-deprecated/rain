package defaults

import (
	"fmt"

	"github.com/wesleimp/rain/pkg/context"
)

// Defaulter intergace
type Defaulter interface {
	fmt.Stringer

	Default(*context.Context) error
}

// Defaulters steps
var Defaulters = []Defaulter{}
