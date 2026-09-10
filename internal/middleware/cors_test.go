package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/fahrigunadi/ytdl/internal/middleware"
)

func setupCORSRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.CORS())
	r.GET("/x", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	return r
}

func TestCORS_SetsHeaders(t *testing.T) {
	r := setupCORSRouter()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected Allow-Origin '*', got %q", got)
	}
	if w.Body.String() != "ok" {
		t.Errorf("handler body should be reached, got %q", w.Body.String())
	}
}

func TestCORS_PreflightShortCircuits(t *testing.T) {
	r := setupCORSRouter()
	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS preflight, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Error("expected Allow-Methods header on preflight response")
	}
}
