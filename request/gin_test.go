package request

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/interviewstreet/logging-go/core"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestData struct {
	namespace string
	uri       string
	method    string
	sortParam string
}

type GinTestSuite struct {
	suite.Suite
	sink   *MemorySink
	router *gin.Engine
	data   TestData
}

func (s *GinTestSuite) SetupSuite() {
	s.data = TestData{
		namespace: "test",
		uri:       "/ping",
		method:    "GET",
		sortParam: "created",
	}
	s.sink = &MemorySink{new(bytes.Buffer)}
	if err := zap.RegisterSink("memory", func(url *url.URL) (zap.Sink, error) {
		return s.sink, nil
	}); err != nil {
		s.Error(err, "Failed to register memory sink with zap")
	}
	gin.SetMode(gin.TestMode)
	s.router = gin.New()
	s.router.Use(GinMiddleware("test", &core.RequestMiddlewareOptions{Env: "test", OutputPath: "memory://"}))
	s.router.Use(gin.Recovery())
}

func (s *GinTestSuite) TestGetRequest() {
	s.router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping?sort=created", nil)
	s.router.ServeHTTP(w, req)

	var output map[string]interface{}
	_ = json.Unmarshal(s.sink.Bytes(), &output)

	s.Equal(s.data.namespace, output["namespace"])
	s.Equal(s.data.uri, output["uri"])
	s.Equal(s.data.method, output["method"])
	urlParams := output["querystring"].(map[string]interface{})
	s.Equal(s.data.sortParam, urlParams["sort"])
}

func TestGinMiddleware(t *testing.T) {
	suite.Run(t, new(GinTestSuite))
}
