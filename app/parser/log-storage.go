package parser

func ParseStorage(str string) map[string]interface{} {
	var details = map[string]interface{}{}

	braceIndex, _ := findNextIndex(str, []string{"{"})

	if braceIndex > 0 {
		jsonText, _ := findJsonText(str[braceIndex:])
		json, _ := convertToJson(jsonText)
		details["connection"] = json
		details["message"] = str[0:braceIndex]
	}else{
		details["message"] = str
	}

	return details
}
