package operators

import (
	"fmt"

	"github.com/maniartech/uexl_go/core"
	"github.com/maniartech/uexl_go/types"
)

func indexer(op string, a, b core.Evaluator, ctx types.Context) (types.Value, error) {
	fmt.Println("Indexer")

	return nil, nil
}

func init() {
	BinaryOpRegistry.Register(".", indexer)
}
