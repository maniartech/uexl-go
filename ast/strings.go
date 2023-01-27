package ast

import (
	"encoding/json"

	"github.com/maniartech/uexl_go/types"
)

type StringNode struct {
	*BaseNode

	Value types.String `json:"value"`
}

// NewStringNode parses the token and returns a StringNode. This function
// currently uses the json.Unmarshal function to parse the string. This
// function is not very efficient and should be replaced with a custom
// function. This function also checks if the string is single quoted or
// double quoted and replaces the first and last character with double
// quotes if the string is single quoted. This is done because the
// json.Unmarshal function only accepts double quoted strings.
func NewStringNode(token []byte, offset, line, col int) (*StringNode, error) {
	var value string
	var singleQuote = false

	// if the first character is a single quote, then it is a single quoted string
	// replace the first and last character with double quotes
	if token[0] == '\'' {
		token[0] = '"'
		token[len(token)-1] = '"'
		singleQuote = true
	}

	err := json.Unmarshal(token, &value)
	if err != nil {
		return nil, err
	}

	// If the string is single quoted, then replace the double quotes with single quotes
	if singleQuote {
		token[0] = '\''
		token[len(token)-1] = '\''
	}

	node := &StringNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeString,
			Line:   line,
			Column: col,
			Offset: offset,
			Token:  string(token),
		},
		Value: types.String(value),
	}

	return node, nil
}

func (n StringNode) Eval(types.Map) (any, error) {
	return n.Value, nil
}
