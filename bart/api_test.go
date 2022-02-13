package bart

import (
	"net/http"
	"testing"
)

type stubAPI struct{ conf *Config }

func (a *stubAPI) clientConf() *Config { return a.conf }

func TestRequestAPI(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		runTest := func(t *testing.T, s stubHandler) {
			server := makeTestServer(t, s)
			defer server.Close()

			client := stubAPI{conf: &Config{HTTP: &http.Client{}}}
			client.conf.baseURL = server.URL

			var out interface{}
			params := apiRequest{route: s.expectedPath, cmd: "foo"}

			got := params.requestAPI(&client, &out)
			if got == nil {
				t.Fatal("expected error, got nil")
			}

			// The error messages seem to change subtly. Rather than test for
			// exact comformance, test that the client can handle all the
			// various shapes of response body. Manually inspect messages here.
			t.Logf(got.Error())
		}

		t.Run("xml", func(t *testing.T) {
			runTest(t, stubHandler{
				expectedPath:     "/err_xml",
				expectedCmd:      "foo",
				responseFilename: "testdata/err.xml",
			})
		})

		t.Run("json object", func(t *testing.T) {
			runTest(t, stubHandler{
				expectedPath:     "/err_json_object",
				expectedCmd:      "foo",
				responseFilename: "testdata/err_object.json",
			})
		})

		t.Run("json string", func(t *testing.T) {
			runTest(t, stubHandler{
				expectedPath:     "/err_json_string",
				expectedCmd:      "foo",
				responseFilename: "testdata/err_string.json",
			})
		})
	})

	t.Run("ok", func(t *testing.T) {
		server := makeTestServer(t, stubHandler{
			expectedPath:     "/ok",
			expectedCmd:      "foo",
			responseFilename: "testdata/ok.json",
		})
		defer server.Close()

		client := stubAPI{conf: &Config{HTTP: &http.Client{}}}
		client.conf.baseURL = server.URL

		// The Root.Message field doesn't always contain an error. Sometimes,
		// successful requests will be an empty string here.
		params := apiRequest{route: "/ok", cmd: "foo"}
		var out interface{}
		got := params.requestAPI(&client, &out)
		if got != nil {
			t.Fatalf("unexpected error, %v", got)
		}
	})
}
