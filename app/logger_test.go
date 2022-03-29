package app

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/interviewstreet/logging-go/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	*bytes.Buffer
}

var sink *MemorySink

// Implement Close and Sync as no-ops to satisfy the interface. The Write
// method is provided by the embedded buffer.
func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

const (
	testTextPayload = "Choosing an Agent"
	testKey         = "agent_name"
	testKeyValue    = "chamber"
	testCtxID       = "match-622"
	namespace       = "test"
)

func setupTest() {
	// Create a sink instance, and register it with zap for the "memory"
	// protocol.
	sink = &MemorySink{new(bytes.Buffer)}
	_ = zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})
	SetupLogger(namespace, zapcore.InfoLevel, &core.LoggerOptions{
		Env: namespace, OutputPath: "memory://", ErrorOutputPath: "memory://",
	})
}

func assertFields(t *testing.T, output, key, value string) {
	compareWith := fmt.Sprintf(`"%s":"%s"`, key, value)
	if !strings.Contains(output, compareWith) {
		t.Errorf("Test failed for `%s` key", key)
	}
}

// TestNew runs test-cases with default generated context-id
func TestNew(t *testing.T) {
	t.Parallel()
	setupTest()

	log := New()
	log.Infow(testTextPayload, testKey, testKeyValue)

	// Assert sink contents
	output := sink.String()
	t.Logf("output = %s", output)

	assertFields(t, output, "text_payload", testTextPayload)
	assertFields(t, output, testKey, testKeyValue)
	assertFields(t, output, "environment", namespace)
	assertFields(t, output, "namespace", namespace)
}

// TestNewWithCtx runs test-cases with provided context-id
func TestNewWithCtx(t *testing.T) {
	t.Parallel()
	setupTest()

	log := NewWithCtx(testCtxID)
	log.Infow(testTextPayload, testKey, testKeyValue)

	// Assert sink contents
	output := sink.String()
	t.Logf("output = %s", output)

	assertFields(t, output, "text_payload", testTextPayload)
	assertFields(t, output, testKey, testKeyValue)
	assertFields(t, output, "environment", namespace)
	assertFields(t, output, "namespace", namespace)
	assertFields(t, output, "context_id", testCtxID)
}

// TestNewWithCtx runs test-cases with provided gin context
func TestNewWithGinCtx(t *testing.T) {
	t.Parallel()
	setupTest()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/select_agent", func(context *gin.Context) {
		logger := NewWithGinCtx(context)
		logger.Infow(testTextPayload, testKey, testKeyValue)
		context.String(200, testKeyValue)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/select_agent", nil)
	req.Header.Set("x-request-id", testCtxID)
	router.ServeHTTP(w, req)

	output := sink.String()
	t.Logf("output = %s", output)

	assertFields(t, output, "text_payload", testTextPayload)
	assertFields(t, output, testKey, testKeyValue)
	assertFields(t, output, "environment", namespace)
	assertFields(t, output, "namespace", namespace)
	assertFields(t, output, "context_id", testCtxID)
}
