package main

import (
"fmt"
"github.com/maniartech/uexl_go/parser"
)

func main() {
expressions := []string{"arr.(i)", "'abc'.(i)"}

for _, expr := range expressions {
fmt.Printf("\nParsing: %s\n", expr)
p := parser.NewParser(expr)
node, err := p.Parse()
if err != nil {
fmt.Printf("Error: %v\n", err)
continue
}

fmt.Printf("Type: %T\n", node)
if idx, ok := node.(*parser.IndexAccess); ok {
fmt.Printf("Index Type: %T\n", idx.Index)
}
}
}
