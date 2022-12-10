package ast

import "strconv"

type Boolean bool

type BooleanNode struct {
	BaseNode

	Value Boolean `json:"value"`
}

func NewBooleanNode(token string, offset, line, col int) (Node, error) {
	value, err := strconv.ParseBool(token)
	if err != nil {
		return nil, err
	}

	node := BooleanNode{
		BaseNode: BaseNode{
			Type:   NodeTypeBoolean,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  token,
		},
		Value: Boolean(value),
	}

	return node, nil
}

func (n BooleanNode) Eval(Map) (any, error) {
	return n.Value, nil
}
