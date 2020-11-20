package parser



func ParseAccess(str string) map[string]interface{} {
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

// Examples
//2020-10-24T13:01:00.035+0000 I  ACCESS   [conn324445] Successfully authenticated as principal __system on local from client 10.9.6.183:51438
//2020-10-24T13:01:05.127+0000 I  ACCESS   [conn324446] Successfully authenticated as principal prod-isl on admin from client 172.26.121.150:56053
//2020-10-24T13:01:05.127+0000 I  ACCESS   [conn324446]  authenticate db: dbprod { authenticate: 1, user: "mbdb", nonce: "xxx", key: "xxx" }
//2020-10-24T13:01:05.127+0000 I  ACCESS   [conn324446] Failed to authenticate mbdb@dbprod with mechanism MONGODB-CR: AuthenticationFailed MONGODB-CR credentials missing in the user document
