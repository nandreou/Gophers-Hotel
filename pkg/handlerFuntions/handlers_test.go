package handlerfuntions

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetHandlers(t *testing.T) {
	mux := setUpRoutes()

	testSrv := httptest.NewTLSServer(mux)
	defer testSrv.Close()

	for _, v := range reqBodyGet {

		response, err := testSrv.Client().Get(testSrv.URL + v.path)

		if err != nil {
			t.Error("Could not Get response from handler")
		} else if v.expected != response.StatusCode {
			t.Error("Could not Get response from handler for the path", v.path)
		}
	}

}

func TestPostHandlers(t *testing.T) {
	mux := setUpRoutes()

	testSrv := httptest.NewTLSServer(mux)
	defer testSrv.Close()

	var data = url.Values{}

	for _, value := range reqBodyPost {

		for _, testValues := range value.keyValue {
			for _, vals := range testValues {
				data.Add(vals.key, vals.value)
			}

			response, err := testSrv.Client().PostForm(testSrv.URL+value.path, data)

			if err != nil {
				t.Error("POST Failed", err)
			} else if response.StatusCode != value.expected[0] && response.StatusCode != value.expected[1] {
				t.Error("POST Failed with", response.StatusCode)
			}

			for keys := range data {
				data.Del(keys)
			}
		}
	}
}
