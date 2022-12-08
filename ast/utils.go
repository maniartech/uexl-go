package ast

func ToNodesSlice(v interface{}) []Node {
	if v == nil {
		return nil
	}
	islice := v.([]interface{})
	// convert iSlice to []Node
	nodes := make([]Node, len(islice))
	for i, node := range islice {
		nodes[i] = node.(Node)
	}
	return nodes
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
