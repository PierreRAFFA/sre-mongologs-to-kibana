package parser

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

func findNextIndex(str string, substrs []string) (int, string) {
	var resultIndex = -1
	var resultSubstr string
	for _, substr := range substrs {
		substrIndex := strings.Index(str, substr)
		if substrIndex != -1 {
			if resultIndex == -1 || substrIndex < resultIndex {
				resultIndex = substrIndex
				resultSubstr = substr
			}
		}
	}
	return resultIndex, resultSubstr
}

func isAlpha(str string) bool {
	const alpha = "abcdefghijklmnopqrstuvwxyz"
	if !strings.Contains(alpha, strings.ToLower(str)) {
		return false
	}
	return true
}

func isNumeric(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

func findJsonText(str string) (string, int) {
	var pos = 0
	var open = 0
	for pos < len(str) {
		character := string(str[pos])
		if character == "{" {
			open = open + 1
		} else if character == "}" {
			open = open - 1
			if open == 0 {
				return str[0 : pos+1], pos + 1
			}
		}
		pos = pos + 1
	}
	return "", 0
}

func convertToJson(jsonString string) (interface{}, error) {
	var result interface{}
	var command string
	re := regexp.MustCompile(` ([.\$\w]+):`)
	command = re.ReplaceAllString(jsonString, ` "$1":`)

	//re = regexp.MustCompile(`UUID\(([-"\w0-9]+)\)`)
	//command = re.ReplaceAllString(command, "$1")

	// Parses json
	err := json.Unmarshal([]byte(command), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

