package bart

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func makeTestServer(t *testing.T, s stubHandler) *httptest.Server {
	return httptest.NewServer(s.makeHandler(t))
}

type stubHandler struct {
	expectedPath     string
	expectedCmd      string
	responseFilename string
}

func (s stubHandler) makeHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		var (
			reqURL *url.URL
			err    error
			file   io.ReadCloser
			data   []byte
		)

		if reqURL, err = url.ParseRequestURI(r.RequestURI); err != nil {
			t.Fatal(err)
		}
		if reqURL.Path != s.expectedPath {
			t.Fatalf("wrong request path; got %q, expected %q", reqURL.Path, s.expectedPath)
		}

		if cmd := reqURL.Query().Get("cmd"); cmd != s.expectedCmd {
			t.Fatalf("unexpected query string value for cmd, %q", cmd)
		}

		if file, err = os.Open(s.responseFilename); err != nil {
			t.Fatal(err)
		}
		defer func() { _ = file.Close() }()
		if data, err = io.ReadAll(file); err != nil {
			t.Fatal(err)
		}

		fmt.Fprintf(w, "%s", data)
	}
}
