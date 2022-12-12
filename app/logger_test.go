package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/interviewstreet/logging-go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	*bytes.Buffer
}

// Implement Close and Sync as no-ops to satisfy the interface. The Write
// method is provided by the embedded buffer.
func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

const sampleHttpURI = "/select_agent"

type TestData struct {
	payload       string
	key           string
	value         string
	contextID     string
	namespace     string
	contextHeader string
}

type TestSuite struct {
	suite.Suite
	sink       *MemorySink
	httpRouter *gin.Engine
	data       TestData
}

func (s *TestSuite) SetupSuite() {
	// Set test data
	s.data = TestData{
		payload:       "Choosing an Agent",
		key:           "agent_name",
		value:         "chamber",
		contextID:     uuid.NewString(),
		namespace:     "test",
		contextHeader: "x-request-id",
	}
	s.sink = &MemorySink{new(bytes.Buffer)}
	if err := zap.RegisterSink("memory", func(url *url.URL) (zap.Sink, error) {
		return s.sink, nil
	}); err != nil {
		s.Error(err, "Failed to register memory sink with zap")
	}
	SetupLogger(s.data.namespace, zapcore.InfoLevel, &core.LoggerOptions{
		Env: s.data.namespace, OutputPath: "memory://",
	})

	// Setup gin
	gin.SetMode(gin.TestMode)
	s.httpRouter = gin.New()
}

func (s *TestSuite) AfterTest(_, _ string) {
	s.sink.Reset()
}

func (s TestSuite) commonTests(output map[string]interface{}) {
	s.Equal(s.data.payload, output["text_payload"], "Payload test failed")
	logContexts := output["labels"].(map[string]interface{})
	s.Equal(s.data.value, logContexts[s.data.key], "Log key-value context test failed")
	s.Equal(s.data.namespace, output["namespace"], "Namespace test failed")
}

func (s *TestSuite) TestNew() {
	log := New()
	log.Infow(s.data.payload, s.data.key, s.data.value)
	var output map[string]interface{}
	_ = json.Unmarshal(s.sink.Bytes(), &output)
	s.commonTests(output)
}

func (s *TestSuite) TestNewWithContextID() {
	log := NewWithContextID(s.data.contextID)
	log.Infow(s.data.payload, s.data.key, s.data.value)
	var output map[string]interface{}
	_ = json.Unmarshal(s.sink.Bytes(), &output)
	s.commonTests(output)
	s.Equal(s.data.contextID, output["context_id"])
}

func (s TestSuite) exampleHTTPRoute(ctx *gin.Context) {
	logger := NewWithGinCtx(ctx)
	logger.Infow(s.data.payload, s.data.key, s.data.value)
	ctx.String(200, s.data.value)
}

func (s *TestSuite) TestNewWithGinCtx() {
	s.httpRouter.GET(sampleHttpURI, s.exampleHTTPRoute)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", sampleHttpURI, nil)
	req.Header.Set(s.data.contextHeader, s.data.contextID)
	s.httpRouter.ServeHTTP(w, req)

	var output map[string]interface{}
	_ = json.Unmarshal(s.sink.Bytes(), &output)
	s.commonTests(output)
	s.Equal(s.data.contextID, output["context_id"])
}

func (s *TestSuite) TestNewWithCtx() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, TraceIDKey, s.data.contextID)

	log := NewWithCtx(ctx)
	log.Infow(s.data.payload, s.data.key, s.data.value)
	var output map[string]interface{}
	_ = json.Unmarshal(s.sink.Bytes(), &output)
	s.commonTests(output)
	s.Equal(s.data.contextID, output["trace_id"])
}

func TestApplicationLoggers(t *testing.T) {
	assert.Panics(t, checkInitialisation)
	suite.Run(t, new(TestSuite))

}
