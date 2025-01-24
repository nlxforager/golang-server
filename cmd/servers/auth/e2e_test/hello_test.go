package e2e_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang-server/cmd/servers/auth/handlers"
)

func TestHandler_Hello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)
	w := httptest.NewRecorder()
	handlers.Hello()(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	_want := "Hello World!"
	if string(data) != _want {
		t.Errorf("expected %s got %v", _want, string(data))
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %v got %v", http.StatusOK, res.StatusCode)
	}
}
