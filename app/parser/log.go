package parser

import (
	"errors"
	"fmt"
	"mongologs/model"
	"regexp"
)

func ParseLog(logString string, id int) (*model.MongoLog, error) {
	var mongolog = &model.MongoLog{}

	var messageString string
	r, _ := regexp.Compile(`^([-+.:T0-9]+) (I|W|E|F|D1|D2|D3|D4|D5)  ([A-Z]+) +\[([-.\w]+)\] (.*)$`)
	matches := r.FindStringSubmatch(logString)
	if len(matches) == 6 {
		mongolog.Id = id
		mongolog.DateTime = matches[1]
		mongolog.Severity = matches[2]
		mongolog.Component = matches[3]
		mongolog.Context = matches[4]
		messageString = matches[5]
	} else {
		return nil, errors.New(fmt.Sprintf("cannot parse `%s`", logString))
	}

	switch mongolog.Component {
	case model.ComponentCommand:
		mongolog.Details = ParseCommand(messageString)
		break

	case model.ComponentNetwork:
		mongolog.Details = ParseNetwork(messageString)
		break

	case model.ComponentAccess:
		mongolog.Details = ParseAccess(messageString)
		break

	case model.ComponentQuery:
		mongolog.Details = ParseQuery(messageString)
		break

	case model.ComponentWrite:
		mongolog.Details = ParseWrite(messageString)
		break

	case model.ComponentStorage:
		mongolog.Details = ParseStorage(messageString)
		break

	case model.ComponentCoonPool:
		mongolog.Details = ParseDefault(messageString)
		break

	default:
		mongolog.Details = ParseDefault(messageString)
	}

	return mongolog, nil
}
