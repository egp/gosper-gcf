// named/e.go v3
package named

import "github.com/egp/gosper-gcf/core"

func E() core.PQStream {
	return mustFiniteRCFPrefixStream([]int64{2, 1, 2, 1, 1, 4, 1, 1, 6})
}

// named/e.go v3
