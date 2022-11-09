package autopilot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

type GraphqlQuery struct {
	Query     string
	Variables map[string]interface{} `json:",omitempty"`
}

func ToJson(query GraphqlQuery) string {
	bytes, _ := json.Marshal(query)
	return string(bytes)
}

func Parse(r *http.Request) GraphqlQuery {
	output := GraphqlQuery{}
	defer r.Body.Close()
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return output
	}
	json.Unmarshal(bytes, &output)
	return output
}

func GraphQLQueryValidation(t *testing.T, exp string) RequestValidation {
	return func(r *http.Request) {
		Equals(t, exp, Parse(r).Query)
	}
}

func GraphQLQueryFixture(fixture string) GraphqlQuery {
	exp := GraphqlQuery{}
	json.Unmarshal([]byte(TemplatedFixture(fixture)), &exp)
	return exp
}

func GraphQLQueryFixtureValidation(t *testing.T, fixture string) RequestValidation {
	return func(r *http.Request) {
		Equals(t, ToJson(GraphQLQueryFixture(fixture)), ToJson(Parse(r)))
	}
}
