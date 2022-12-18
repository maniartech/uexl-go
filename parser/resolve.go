package parser

func resolveAscii(bytes interface{}) string {
	var pType string
	for _, eachByte := range bytes.([]interface{}) {
		for _, byArr := range eachByte.([]interface{}) {
			pType += string(byArr.([]byte))
		}
	}

	return pType
}
