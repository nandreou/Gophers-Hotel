package handlerfuntions

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

	for _, value := range reqAvailaibilityBodyPost {
		data.Add("start", value.start)
		data.Add("end", value.end)

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

	for _, value := range reqBodyPost {

		for _, testValues := range value.keyValue {
			for _, vals := range testValues {
				data.Add(vals.key, vals.value)
			}

			encodedData := data.Encode()

			r, err := http.NewRequest("POST", "/book-now/1?ri=test", strings.NewReader(encodedData))

			if err != nil {
				t.Error("POST Failed", err)
			}

			ctx := MakeNewCtx(r)
			r = r.WithContext(ctx)

			rr := httptest.NewRecorder()

			Repo.App.Session.Put(ctx, "test", "something")

			handler := http.HandlerFunc(Repo.BookNow)

			handler.ServeHTTP(rr, r)
			fmt.Println(rr)

			for keys := range data {
				data.Del(keys)
			}
		}
	}
}

func MakeNewCtx(r *http.Request) context.Context {
	ctx, err := app.Session.Load(r.Context(), r.Header.Get("X-Session"))

	if err != nil {
		log.Println(err)
		return nil
	}
	return ctx
}
