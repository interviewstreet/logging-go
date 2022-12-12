package core

type LoggerOptions struct {
	Env             string `default:"production"`
	OutputPath      string `default:"stdout"`
	ErrorOutputPath string `default:"stderr"`
	ContextIDHeader string `default:"x-request-id"`
}

type RequestMiddlewareOptions struct {
	Env             string   `default:"production"`
	OutputPath      string   `default:"stdout"`
	ContextIDHeader string   `default:"x-request-id"` // The header key in which the unique context-id is expected
	IgnoredPaths    []string // List of relative paths (without host) to be ignored from logging
	AddTraceID      bool     `default:"false"` // Flag to control whether trace_id should be added to the request's context
}
