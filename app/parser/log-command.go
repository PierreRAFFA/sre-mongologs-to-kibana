package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParseCommand(str string) map[string]interface{} {
	var details = map[string]interface{}{}
	r, _ := regexp.Compile(`^\w+ ([$.\w]+) (?:appName: "(.*)" )?command:(?: ([\w]+))? (.*) ([0-9]+)ms$`)
	matches := r.FindStringSubmatch(str)

	var command string
	if len(matches) < 6 {
		details["message"] = str
		return details
	}

	details["namespace"] = matches[1]
	details["appName"] = matches[2]
	details["commandName"] = matches[3]

	matches5, _ := strconv.Atoi(matches[5])
	details["duration"] = matches5
	command = matches[4]

	// Extract parts from the command
	var pos = 0
	var currentField = "command"
	for pos < len(command) {
		character := string(command[pos])

		if character == "{" {
			str, nextpos := findJsonText(command[pos:])
			fmt.Sprintf(str)

			if currentField == "command" || currentField == "originatingCommand" {
				details[currentField] = str
			} else {
				json, _ := convertToJson(str)
				details[currentField] = json
			}

			pos = pos + nextpos
		} else if isAlpha(character) {
			charIndex, charDelimiter := findNextIndex(command[pos:], []string{":", " "})
			if charDelimiter == ":" {
				currentField = command[pos : pos+charIndex]
				pos = pos + charIndex + 1
			} else if charDelimiter == " " {
				details[currentField] = command[pos : pos+charIndex]
				pos = pos + charIndex + 1

				// if it was planSummary, we potentially expect planSummaryDetails straight after
				if currentField == "planSummary" {
					currentField = "planSummaryDetails"
				}
			} else {
				// end of command here
				details[currentField] = command[pos:]
				pos = len(command)
			}
		} else if isNumeric(character) {
			spaceIndex := strings.Index(command[pos:], " ")
			if spaceIndex >= 0 {
				value := command[pos : pos+spaceIndex]

				// Here, we force the field to be string because it's a too high number value
				// And makes ES struggle to filter by that value
				if currentField == "cursorid" {
					details[currentField] = value
				}else{
					numericValue, _ := strconv.Atoi(value)
					details[currentField] = numericValue
				}

				pos = pos + spaceIndex + 1
			} else {
				value := command[pos:]

				// Here, we force the field to be string because it's a too high number value
				// And makes ES struggle to filter by that value
				if currentField == "cursorid" {
					details[currentField] = value
				}else{
					numericValue, _ := strconv.Atoi(value)
					details[currentField] = numericValue
				}

				pos = len(command)
			}
		} else {
			// probably a space
			pos = pos + 1
		}
	}

	return details
}