package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Scalingo/sclng-backend-test-v1/internal/handlers"
)

func TestHealthHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		expected   string
		statusCode int
	}{
		{
			name:       "success",
			method:     http.MethodGet,
			expected:   "OK",
			statusCode: http.StatusOK,
		},
		{
			name:       "failure with wrong method",
			method:     http.MethodPost,
			expected:   "only GET method is allowed\n",
			statusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/health", nil)
			w := httptest.NewRecorder()

			handlers.HealthHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
			}

			if tc.statusCode != res.StatusCode {
				t.Errorf("expected %d, got %d", tc.statusCode, res.StatusCode)
			}

			if tc.expected != string(data) {
				t.Errorf("expected %s,  got %s", tc.expected, string(data))
			}
		})
	}
}
