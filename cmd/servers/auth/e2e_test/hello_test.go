package e2e_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"net/http/httptest"
	"testing"

	auth "golang-server/cmd/servers/auth/mux"
)

func TestHandler_Hello(t *testing.T) {
	req := httptest.NewRequestWithContext(context.TODO(), http.MethodGet, "/", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	auth.NewMux().ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
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

	_want := "helloworld"
	if resp.Data.Message != _want {
		t.Errorf("expected %s got %v", _want, resp.Data.Message)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}
}
