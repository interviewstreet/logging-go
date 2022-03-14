package request

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/interviewstreet/logging-go/core"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func setupRouter() *gin.Engine {
	router := gin.New()
	router.Use(GinMiddleware("test", &core.RequestMiddlewareOptions{Env: "test", OutputPath: "memory://"}))
	router.Use(gin.Recovery())
	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})
	return router
}

func TestGinMiddleware(t *testing.T) {
	// Create a sink instance, and register it with zap for the "memory"
	// protocol.
	sink := &MemorySink{new(bytes.Buffer)}
	_ = zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping?sort=created", nil)
	router.ServeHTTP(w, req)

	// Assert sink contents
	output := sink.String()
	t.Logf("output = %s", output)

	assertFields(t, output, "namespace", "test")
	assertFields(t, output, "uri", "/ping")
	assertFields(t, output, "environment", "test")
	assertFields(t, output, "method", "GET")
	assertFields(t, output, "sort", "created")
}
