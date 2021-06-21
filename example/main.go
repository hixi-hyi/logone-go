package main

import (
	"context"
	"strconv"

	"github.com/hixi-hyi/logone-go/logone"
	"github.com/pkg/errors"
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

	if err := doErr(); err != nil {
		//logger.Error("doErr() fmt: %+v", err)
		//logger.Error("doErr() withError").WithError(err)
	}
	if err := doWithStackErr(); err != nil {
		logger.Error("doWithStackErr() fmt: %+v", err)
		logger.Error("doWithStackErr() withError").WithError(err)
		err = errors.WithStack(err)
		logger.Error("doWithStackErr() fmt: %+v", err)
		logger.Error("doWithStackErr() withError").WithError(err)
	}

	if err := doWrapErr(); err != nil {
		err = errors.Wrap(err, "wrap2")
		//logger.Error("doWrapErr() fmt: %+v", err)
		//logger.Error("doWrapErr() withError").WithError(err)
	}

	return
}

func doErr() error {
	if _, err := strconv.ParseInt("", 32, 10); err != nil {
		return err
	}
	return nil
}
func doWithStackErr() error {
	if _, err := strconv.ParseInt("", 32, 10); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
func doWrapErr() error {
	if _, err := strconv.ParseInt("", 32, 10); err != nil {
		err = errors.Wrap(err, "wrap")
		return err
	}
	return nil
}
