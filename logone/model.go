package logone

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

type LogEntry struct {
	Severity   Severity    `json:"severity"`
	Message    string      `json:"message"`
	Time       time.Time   `json:"time,omitempty"`
	Filename   string      `json:"filename,omitempty"`
	Fileline   int         `json:"fileline,omitempty"`
	Funcname   string      `json:"funcname,omitempty"`
	Tags       []string    `json:"tags,omitempty"`
	Elapsed    float64     `json:"elapsed,omitempty"`
	Attributes interface{} `json:"attributes,omitempty"`
	Error      string      `json:"error,omitempty"`
	StackTrace []string    `json:"stackTrace,omitempty"`
}

func (lr *LogEntry) WithTags(tags ...string) *LogEntry {
	lr.Tags = tags
	return lr
}

func (lr *LogEntry) WithAttributes(i interface{}) *LogEntry {
	lr.Attributes = i
	return lr
}

func (lr *LogEntry) WithError(err error) *LogEntry {
	if err == nil {
		return lr
	}
	lr.Error = err.Error()

	// For errors.WithStack
	msg := fmt.Sprintf("%+v", err)
	msg = strings.Replace(msg, "\n\t", " - ", -1)
	lr.StackTrace = strings.Split(msg, "\n")[1:];
	for i := 0; i < len(lr.StackTrace); i++ {
		if strings.HasPrefix(lr.StackTrace[i], "runtime.goexit") {
			lr.StackTrace = lr.StackTrace[0:i+1]
			break;
		}
	}
	return lr
}

type LogRequest struct {
	Type    string      `json:"type"`
	Context *LogContext `json:"context"`
	Runtime *LogRuntime `json:"runtime"`
	Config  *LogConfig  `json:"config"`
	//Extras  interface{}       `json:"extras,omitempty"`
}

type LogConfig struct {
	ElapsedUnit string `json:"elapsedUnit"`
}

type LogRuntime struct {
	Severity   Severity      `json:"severity"`
	StartTime  time.Time     `json:"startTime"`
	EndTime    time.Time     `json:"endTime"`
	Elapsed    int64         `json:"elapsed"`
	Lines      []*LogEntry   `json:"lines,omitempty"`
	Tags       LogTags       `json:"tags,omitempty"`
	Severities SeverityCount `json:"-"`
}

func NewLogRuntime() *LogRuntime {
	lr := &LogRuntime{}
	lr.Tags = LogTags{}
	lr.Severities = SeverityCount{}
	return lr
}

func (lr *LogRuntime) AppendLogEntry(l *LogEntry) {
	lr.Lines = append(lr.Lines, l)
}

type LogContext map[string]string

func NewLogContext() *LogContext {
	lc := &LogContext{}
	return lc
}

type LogTags map[string]int64

func (lt LogTags) CountUp(tags ...string) {
	for _, t := range tags {
		lt[t] += 1
	}
}

type Severity int

const (
	SeverityUnknown Severity = iota
	SeverityDebug
	SeverityInfo
	SeverityWarning
	SeverityError
	SeverityCritical
)

var Severities = map[Severity]string{
	SeverityDebug:    "DEBUG",
	SeverityInfo:     "INFO",
	SeverityWarning:  "WARNING",
	SeverityError:    "ERROR",
	SeverityCritical: "CRITICAL",
}

func (s Severity) String() string {
	if v, ok := Severities[s]; ok {
		return v
	}
	return "UNKNOWN"
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type SeverityCount map[Severity]int64

func (sc SeverityCount) CountUp(s Severity) {
	sc[s]++
}
func (sc SeverityCount) HighestSeverity() Severity {
	keys := []int{int(SeverityUnknown)}
	for k := range sc {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return Severity(keys[len(keys)-1])
}
