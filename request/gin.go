package request

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mcuadros/go-defaults"

	"github.com/interviewstreet/logging-go/app"
	"github.com/interviewstreet/logging-go/core"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// getQueryParams parses request for query parameters and converts them to a map
func getQueryParams(req *http.Request) map[string]string {
	params := req.URL.Query()
	result := make(map[string]string)
	for key, value := range params {
		result[key] = value[0]
	}
	return result
}

// getUrl returns the complete URL for the request
func getUrl(req *http.Request) string {
	return fmt.Sprintf("%s%s", req.Host, req.URL.Path)
}

// getUri returns the resource path for a request
func getUri(req *http.Request) string {
	return req.URL.Path
}

// getContextId extracts the unique id for the request that has been sent, or else generates a new one and
// assigns the value to the request header
func getContextId(req *http.Request, header string) string {
	val := req.Header.Get(header)
	if val == "" {
		val = uuid.New().String()
		req.Header.Set(header, val)
	}
	return val
}

// cleanHeaders parses the request headers and converts them to map
func cleanHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, val := range headers {
		result[key] = val[0]
	}
	return result
}

// GinMiddleware defines the gin handler to intercept all requests and logging the information in a structured manner.
func GinMiddleware(namespace string, options *core.RequestMiddlewareOptions) gin.HandlerFunc {
	// Ensure logger is initialised only once
	loggerOnce.Do(func() {
		if options == nil {
			options = &core.RequestMiddlewareOptions{}
		}
		defaults.SetDefaults(options)
		baseLogger = createNewLogger(namespace, &core.LoggerOptions{Env: options.Env, OutputPath: options.OutputPath})
	})

	return func(c *gin.Context) {
		uri := getUri(c.Request)

		// Don't log if path is ignored
		if contains(options.IgnoredPaths, uri) {
			return
		}

		t0 := time.Now()

		if options.AddTraceID {
			ctx := c.Request.Context()
			ctx = context.WithValue(ctx, app.TraceIDKey, getContextId(c.Request, options.ContextIDHeader))
			c.Request = c.Request.WithContext(ctx)
		}

		// Extract information before request execution
		fields := []interface{}{
			"client_ip", c.ClientIP(),
			"method", c.Request.Method,
			"request_headers", cleanHeaders(c.Request.Header),
			"url", getUrl(c.Request),
			"uri", uri,
			"querystring", getQueryParams(c.Request),
			"trace_id", getContextId(c.Request, options.ContextIDHeader),
		}

		// Wait for the request controller to return
		c.Next()

		fields = append(fields, "latency", time.Since(t0).Microseconds())
		response := c.Writer
		fields = append(fields, "status", response.Status(), "response_headers", cleanHeaders(response.Header()))
		baseLogger.Infow("", fields...)
	}
}
