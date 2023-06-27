package logger

import "github.com/hashicorp/go-hclog"

type Logger struct {
	hclog.Logger
}

func NewLogger() Logger {
	return Logger{
		newLogger(),
	}
}

func newLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name: "auth-service",
	})
}
