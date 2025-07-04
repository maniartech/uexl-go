package ast

import (
	"encoding/json"

	"github.com/maniartech/uexl_go/types"
)

type StringNode struct {
	*BaseNode
	Value          types.String `json:"value"`
	IsSingleQuoted bool         `json:"isSingleQuoted,omitempty"`
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
	var isSingleQuoted bool
	var tokenStr string

	if token[0] == 'r' {
		// Raw string: r"..." or r'...'
		// Check if it's single-quoted or double-quoted
		if len(token) > 2 && token[1] == '\'' {
			isSingleQuoted = true
		} else {
			isSingleQuoted = false
		}
		value = string(token[2 : len(token)-1])
		tokenStr = string(token)
	} else if token[0] == '\'' {
		// Single-quoted string: manually extract content between quotes
		content := string(token[1 : len(token)-1])
		// For single-quoted strings, we process escape sequences manually
		// This avoids issues with json.Unmarshal when the string contains unescaped double quotes
		value = content
		isSingleQuoted = true
		tokenStr = string(token)
	} else {
		// Double-quoted string
		err := json.Unmarshal(token, &value)
		if err != nil {
			return nil, err
		}
		isSingleQuoted = false
		tokenStr = string(token)
	}

	node := &StringNode{
		BaseNode: &BaseNode{
			Type:   NodeTypeString,
			Line:   line,
			Column: col,
			Token:  tokenStr,
		},
		Value:          types.String(value),
		IsSingleQuoted: isSingleQuoted,
	}
	return node, nil
}

func (n StringNode) Eval(types.Context) (types.Value, error) {
	return n.Value, nil
}
