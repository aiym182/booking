package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type PostData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []PostData
	expectedStatusCode int
}{
	{"home", "/", "GET", []PostData{}, http.StatusOK},
	{"about", "/about", "GET", []PostData{}, http.StatusOK},
	{"gq", "/generals-quaters", "GET", []PostData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []PostData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []PostData{}, http.StatusOK},
	{"contact", "/contact", "GET", []PostData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []PostData{}, http.StatusOK},
	{"post-search-avail", "/search-availability", "POST", []PostData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"post-search-avail-json", "/search-availability-json", "POST", []PostData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-02"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []PostData{
		{key: "first_name", value: "john"},
		{key: "last_name", value: "smith"},
		{key: "email", value: "aiym182@gmail.com"},
		{key: "phone", value: "432242323"},
	}, http.StatusOK},
}

func TestHand(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)

	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Error(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Error(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
