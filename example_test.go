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

// Request:
// curl --location --request GET 'http://localhost:8080/add'
// Output:
// {"severity":"INFO","timestamp":1647258527608.217,"source_caller":"logging-go/example.go:32","source_function":"main.AddController","text_payload":"Adding 2 & 5","environment":"dev","namespace":"example","resource_labels":{"hostname":"Sohels-MacBook-Pro-2.local","pid":47460},"context_id":"3e4547da-61d0-44c9-87bd-d88ba2c80298","labels":{}}
// {"timestamp":1647258527608.607,"environment":"test","namespace":"test","resource_labels":{"hostname":"Sohels-MacBook-Pro-2.local","pid":47460},"client_ip":"::1","method":"GET","request_headers":{"Accept":"*/*","Accept-Encoding":"gzip, deflate, br","Connection":"keep-alive","Postman-Token":"9a8e1488-bfdf-46d0-977e-5552f67fefa6","User-Agent":"PostmanRuntime/7.29.0"},"url":"localhost:8080/add","uri":"/add","querystring":{},"context_id":"3e4547da-61d0-44c9-87bd-d88ba2c80298","latency":2006,"status":200,"response_headers":{"Content-Type":"text/plain; charset=utf-8"}}
//
// Request:
// curl --location --request GET 'http://localhost:8080/subtract' --header 'X-Request-ID: 123456'
// Output:
// {"severity":"INFO","timestamp":1647258574408.834,"source_caller":"logging-go/example.go:43","source_function":"main.SubtractController","text_payload":"Subtract 2 from 5","environment":"dev","namespace":"example","resource_labels":{"hostname":"Sohels-MacBook-Pro-2.local","pid":47460},"context_id":"123456","labels":{}}
// {"timestamp":1647258574408.907,"environment":"test","namespace":"test","resource_labels":{"hostname":"Sohels-MacBook-Pro-2.local","pid":47460},"client_ip":"::1","method":"GET","request_headers":{"Accept":"*/*","Accept-Encoding":"gzip, deflate, br","Connection":"keep-alive","Postman-Token":"53069711-6dc3-4c6b-a144-c7f8c3704b3f","User-Agent":"PostmanRuntime/7.29.0","X-Request-Id":"123456"},"url":"localhost:8080/subtract","uri":"/subtract","querystring":{},"context_id":"123456","latency":103,"status":200,"response_headers":{"Content-Type":"text/plain; charset=utf-8"}}
