package parser

func resolveAscii(bytes interface{}) string {
	var pType string
	for _, eachByte := range bytes.([]interface{}) {
		pType += string(eachByte.([]byte))
	}

	return pType
}
