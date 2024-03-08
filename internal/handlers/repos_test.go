package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Scalingo/sclng-backend-test-v1/internal/handlers"
)

func TestReposHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/repos", nil)
	w := httptest.NewRecorder()

	handlers.ReposHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("failed to read body: %v", err)
	}

	fmt.Println(data)
}
