package core

type LoggerOptions struct {
	Env             string `default:"production"`
	OutputPath      string `default:"stdout"`
	ErrorOutputPath string `default:"stderr"`
	ContextIDHeader string `default:"x-request-id"`
}

type RequestMiddlewareOptions struct {
	Env             string `default:"production"`
	OutputPath      string `default:"stdout"`
	ContextIDHeader string `default:"x-request-id"` // The header key in which the unique context-id is expected
}
