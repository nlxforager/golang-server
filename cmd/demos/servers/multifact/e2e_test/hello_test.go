package e2e_test

import (
	"context"
	"encoding/json"
	"golang.org/x/net/html"
	"io"
	"net/http"

	"net/http/httptest"
	"testing"

	auth "golang-server/cmd/demos/servers/multifact/mux"
)

func TestHandler_Hello_AcceptJSON(t *testing.T) {
	req := httptest.NewRequestWithContext(context.TODO(), http.MethodGet, "/", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	auth.NewMux(nil).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	type response struct {
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	}

	var resp response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v, data %s", err, string(data))
	}

	_want := "Hello World"
	if resp.Data.Message != _want {
		t.Errorf("expected %s got %v", _want, resp.Data.Message)
	}

}

func TestHandler_Hello_AcceptUnspecified(t *testing.T) {
	req := httptest.NewRequestWithContext(context.TODO(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Accept", "aaa/json")

	auth.NewMux(nil).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotAcceptable {
		t.Errorf("expected status code to be %v got %v", http.StatusNotAcceptable, res.StatusCode)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	type response struct {
		Data struct {
			Message string `json:"message"`
		} `json:"data"`
	}

	var resp response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	_want := ""
	if resp.Data.Message != _want {
		t.Errorf("expected %s got %v", _want, resp.Data.Message)
	}
}

func TestHandler_Hello_AcceptHTML(t *testing.T) {
	req := httptest.NewRequestWithContext(context.TODO(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	req.Header.Set("Accept", "text/html")

	auth.NewMux(nil).ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}

	data, err := html.Parse(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	//t.Logf("??  HTML %#v", data.FirstChild)
	//t.Logf("??  HEAD = HTML.child %#v", data.FirstChild.FirstChild)
	//t.Logf("??  BODY = HEAD.Next %#v", data.FirstChild.FirstChild.NextSibling)
	//t.Logf("??  div = BODY.child %#v", data.FirstChild.FirstChild.NextSibling.FirstChild)
	want := data.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.Data
	//t.Logf("??  WANT = div.child.data %#v", data.FirstChild.FirstChild.NextSibling.FirstChild.FirstChild.Data)

	if data.FirstChild.Data != "html" {
		t.Errorf("expected %s got %#v", "html", data.FirstChild.Data)
	}
	if want != "Hello World" {
		t.Errorf("expected %s got %#v", "want", want)
	}

}
