package ast

import (
	"bytes"
	"fmt"
)

type ObjectItem struct {
	Key   String
	Value Node
}

type ObjectNode struct {
	BaseNode

	Items []ObjectItem
}

func NewObjectNode(token string, items []ObjectItem, offset, line, col int) (Node, error) {
	node := ObjectNode{
		BaseNode: BaseNode{
			Type:   NodeTypeObject,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},

		Items: items,
	}

	return node, nil
}

func (n ObjectNode) String() string {
	var b bytes.Buffer
	b.WriteString("{")
	for i, item := range n.Items {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%s: %v", item.Key, item.Value))
	}
	b.WriteString("}")

	return fmt.Sprintf("ObjectNode: %s", b.String())
}

// Eval parses the array node and returns the array node.
func (n ObjectNode) Eval(env Map) (any, error) {
	results := make(map[String]interface{})
	for _, item := range n.Items {
		result, err := item.Value.Eval(env)
		if err != nil {
			return nil, err
		}
		results[item.Key] = result
	}
	return results, nil
}
