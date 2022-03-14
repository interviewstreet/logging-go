# Go Logging SDK
Go SDK for supporting structured logging in applications

## Usage
```go
package logginggo_test

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/interviewstreet/logging-go/app"
	"github.com/interviewstreet/logging-go/core"
	"github.com/interviewstreet/logging-go/request"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	// Setup the application logger once
	app.SetupLogger("example", zapcore.InfoLevel, &core.LoggerOptions{Env: "dev"})
}

func Add(a, b int, log *zap.SugaredLogger) int {
	log.Debugw("Adding two numbers", "a", a, "b", b)
	return a + b
}

func Subtract(a, b int, log *zap.SugaredLogger) int {
	log.Debugw("Subtracting two numbers", "a", a, "b", b)
	return a - b
}

func AddController(ctx *gin.Context) {
	// Get a new logger object
	log := app.New()
	log.Info("Adding 2 & 5")

	// Pass the parent logger to child functions
	result := Add(2, 5, log)

	ctx.String(200, strconv.Itoa(result))
}

func SubtractController(ctx *gin.Context) {
	// Get a new logger object with pre-existing context
	log := app.NewWithCtx(ctx.Request.Header.Get("x-request-id"))
	log.Info("Subtract 2 from 5")

	// Pass the parent logger to child functions
	result := Subtract(5, 2, log)

	ctx.String(200, strconv.Itoa(result))
}

func Example() {
	router := gin.New()
	// Add the request logging middleware
	router.Use(
		request.GinMiddleware("test", &core.RequestMiddlewareOptions{Env: "test"}),
	)
	router.Use(gin.Recovery())
	router.GET("/add", AddController)
	router.GET("/subtract", SubtractController)
	router.Run()
}
```