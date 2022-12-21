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

func resolveAscii1(bytes interface{}) string {
	var val string
	for _, eachByte := range bytes.([]interface{}) {
		val += string(eachByte.([]byte))
	}

	return val
}
