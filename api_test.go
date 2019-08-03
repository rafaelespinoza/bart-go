package bart

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRequestAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqURL, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			panic(err)
		}
		switch reqURL.Path {
		case "/route.aspx":
			fmt.Fprintf(w, `<?xml version="1.0" encoding="utf-8"?>
<root>
	<message>
		<error>
			<text>Invalid cmd</text>
			<details>The cmd parameter (%s) is missing or invalid. Please correct the error and try again.</details>
		</error>
	</message>
</root>`, reqURL.Query().Get("cmd"))
		default:
			err := fmt.Errorf("unknown request path %s", reqURL.Path)
			panic(err)
		}
	}))
	originalBaseURL := baseURL
	baseURL = server.URL
	defer func() {
		server.Close()
		baseURL = originalBaseURL
	}()

	var out interface{}
	actual := requestAPI("/route.aspx", "foo", nil, &out)
	expectedError := "error: Invalid cmd. The cmd parameter (foo) is missing or invalid. Please correct the error and try again."
	if actual.Error() != expectedError {
		t.Errorf("Error formatted incorrectly. Got %q, expected %q", actual.Error(), expectedError)
	}
}
