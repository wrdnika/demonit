package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireAPIKey_AcceptsValidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/secure", requireAPIKey(HeaderDeviceAPIKey, "secret-key"), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set(HeaderDeviceAPIKey, "secret-key")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRequireAPIKey_RejectsMissingOrInvalidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/secure", requireAPIKey(HeaderDeviceAPIKey, "secret-key"), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	cases := []struct {
		name   string
		header string
	}{
		{name: "missing", header: ""},
		{name: "wrong", header: "nope"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/secure", nil)
			if tc.header != "" {
				req.Header.Set(HeaderDeviceAPIKey, tc.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("expected 401, got %d", w.Code)
			}

			var body APIError
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body.Error.Code != "UNAUTHORIZED" {
				t.Fatalf("expected UNAUTHORIZED code, got %q", body.Error.Code)
			}
		})
	}
}
