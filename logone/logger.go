package logone

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type Logger struct {
	Config     *Config
	LogRequest *LogRequest
}

func NewLogger(c *Config) *Logger {
	lr := &LogRequest{
		Type:    c.Type,
		Context: NewLogContext(),
		Runtime: NewLogRuntime(),
		Config: &LogConfig{
			ElapsedUnit: c.ElapsedUnit.String(),
		},
	}

	l := &Logger{
		Config:     c,
		LogRequest: lr,
	}
	return l
}

func NewLoggerDefault() *Logger {
	c := NewConfigDefault()
	return NewLogger(c)
}

func (l *Logger) SetLogContext(lc *LogContext) {
	l.LogRequest.Context = lc
}

func (l *Logger) Start() func() {
	lr := l.LogRequest.Runtime
	lr.StartTime = time.Now()
	return func() {
		l.Finish()
	}
}

func (l *Logger) FillInRuntimeMetadata() {
	lr := l.LogRequest.Runtime

	lr.EndTime = time.Now()
	elapsed := lr.EndTime.Sub(lr.StartTime)
	lr.Elapsed = int64(elapsed / l.Config.ElapsedUnit)

	for _, line := range lr.Lines {
		lr.Tags.CountUp(line.Tags...)
		lr.Severities.CountUp(line.Severity)
	}
	lr.Severity = lr.Severities.HighestSeverity()
}

func (l *Logger) Finish() {
	l.FillInRuntimeMetadata()
	var logline []byte
	if l.Config.JsonIndent {
		logline, _ = json.MarshalIndent(l.LogRequest, "", "  ")
	} else {
		logline, _ = json.Marshal(l.LogRequest)
	}
	fmt.Println(string(logline))
}

func (l *Logger) Record(severity Severity, message string) *LogEntry {
	funcname, _, fileline := FileInfo(3)
	e := &LogEntry{
		Severity: severity,
		Message:  message,
		Time:     time.Now(),
		//Filename: filename,
		//It is commented out because the filename of the build environment is displayed.
		Fileline: fileline,
		Funcname: funcname,
	}
	l.LogRequest.Runtime.AppendLogEntry(e)
	return e
}

func (l *Logger) Debug(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityDebug, fmt.Sprintf(f, v...))
}
func (l *Logger) Info(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityInfo, fmt.Sprintf(f, v...))
}
func (l *Logger) Warning(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityWarning, fmt.Sprintf(f, v...))
}
func (l *Logger) Error(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityError, fmt.Sprintf(f, v...))
}
func (l *Logger) Critical(f string, v ...interface{}) *LogEntry {
	return l.Record(SeverityCritical, fmt.Sprintf(f, v...))
}

func FileInfo(depth int) (string, string, int) {
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		return "???", "???", 0
	}
	fn := runtime.FuncForPC(pc)
	filename, fileline := fn.FileLine(pc)
	return fn.Name(), filename, fileline
}
