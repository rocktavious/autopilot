package autopilot

import (
	"encoding/json"
)

type TestRequest struct {
	Request  GraphqlQuery
	Response map[string]any
}

func NewTestRequest(request string, variables string, response string) TestRequest {
	testRequest := TestRequest{
		Request: GraphqlQuery{
			Query:     request,
			Variables: templatedJson(variables),
		},
		Response: templatedJson(response),
	}
	return testRequest
}

func templatedJson(values string) map[string]any {
	parsedValues, err := Templater.Use(values)
	if err != nil {
		panic(err)
	}
	var valuesJSON map[string]any
	if err := json.Unmarshal([]byte(parsedValues), &valuesJSON); err != nil {
		panic(err)
	}

	return valuesJSON
}
