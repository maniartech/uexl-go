package parser

func ParseString(input string) (Node, error) {
	return NewParser(input).Parse()
}
