package logone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	t.Run("scenario", func(t *testing.T) {
		l := NewLoggerDefault()
		l.SetLogContext(&LogContext{
			"REQUEST_ID": "xxxx",
		})
		finish := l.Start()
		attribute := map[string]string{"test": "1"}
		l.Debug("%s", "debug").WithTags("aws-sdk-error", "notify")
		l.Critical("%s", "critical").WithTags("aws-sdk-error", "error").WithAttributes(attribute)
		finish()

		assert.Exactly(t, int64(2), l.LogRequest.Runtime.Tags["aws-sdk-error"])
		assert.Exactly(t, int64(1), l.LogRequest.Runtime.Tags["notify"])
		assert.Exactly(t, int64(1), l.LogRequest.Runtime.Tags["error"])
		assert.Exactly(t, SeverityCritical, l.LogRequest.Runtime.Severity)

		assert.Exactly(t, "xxxx", (*l.LogRequest.Context)["REQUEST_ID"])
		// TODO test to output line
	})

}
