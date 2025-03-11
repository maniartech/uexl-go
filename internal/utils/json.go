package utils

import "encoding/json"

func ToJSONString(v interface{}) string {
	jsonBytes, _ := json.MarshalIndent(v, "", "  ")
	return string(jsonBytes)
}

func PrintJSON(v interface{}) {
	println(ToJSONString(v))
}
