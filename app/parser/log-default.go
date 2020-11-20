package parser

func ParseDefault(str string) map[string]interface{} {
	var details = map[string]interface{}{}
	details["message"] = str
	return details
}