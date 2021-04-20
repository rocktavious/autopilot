package "autopilot"

import (
	"fmt"
  "bytes"
	"io/ioutil"
	"net/http"
  "net/http/httptest"
	"path/filepath"
	"runtime"
	"reflect"
	"testing"
)

/*
var (
  Client *MyClient
)

func TestMain(m *testing.M) {
	flag.Parse()
  teardown := autopilot.Setup()
  defer teardown()
  Client := NewClient(autopilot.Server.URL)
	os.Exit(m.Run())
}

func TestExample(t *testing.T) {
  // Arrange
	autopilot.Mux.HandleFunc("/orgs/octokit/repos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, autopilot.Fixture("repos/octokit.json"))
	})
  // Act
	body, err := Client.DoStuff()
  a := Client.GetThing("Four")
  // Assert
  autopilot.Assert(t, a >= 4, "value must be greater than or equal to four")
  autopilot.Ok(t, err)
	autopilot.Equals(t, []byte("OK"), body)
}
*/

var (
	Mux    *http.ServeMux
	Server *httptest.Server
)

// Setup an HttpTestServer and ServerMux (for path routing) and return the teardown function
func Setup() func() {
	Mux = http.NewServeMux()
	Server = httptest.NewServer(Mux)
	return func() {
		Server.Close()
	}
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// Returns a Mock'd HttpClient
func NewHttpClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
/*
	client := autopilot.NewHttpClient(func(req *http.Request) *http.Response {
		// Test request parameters
		autopilot.Equals(t, req.URL.String(), "http://example.com/some/path")
		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
 			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	api := API{client, "http://example.com"}
*/

// Load testdata/fixtures/<path> and return as string
func Fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
