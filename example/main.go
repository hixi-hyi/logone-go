package main

import (
	"context"

	"github.com/hixi-hyi/logone-go/logone"
)

func main() {
	ctx := context.Background()
	manager := logone.NewManagerDefault()
	ctx, finish := manager.RecordingWithContext(ctx)
	defer finish()
	logger, _ := logone.LoggerFromContext(ctx)
	logger.SetLogContext(&logone.LogContext{
		"REQUEST_ID": "xxxxxxx",
	})
	ctx = logone.NewContextWithLogger(ctx, logger)
	logger.Debug("invoked").WithTags("critical").WithAttributes("xxxx")

	return
}
