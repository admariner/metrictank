package expr

import (
	"sync"

	"github.com/grafana/metrictank/schema"
)

func init() {
	pointSlicePool = &sync.Pool{
		New: func() interface{} { return make([]schema.Point, 0, 100) },
	}
}
