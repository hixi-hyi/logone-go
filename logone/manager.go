package logone

import (
	"context"
	"time"
)

type Config struct {
	Type            string
	DefaultSeverity Severity
	OutputSeverity  Severity
	JsonIndent      bool
	ElapsedUnit     time.Duration
}

type Manager struct {
	Config *Config
}

func NewConfigDefault() *Config {
	return &Config{
		Type:            "request",
		DefaultSeverity: SeverityDebug,
		OutputSeverity:  SeverityDebug,
		ElapsedUnit:     time.Millisecond,
	}
}

func NewManager(mc *Config) *Manager {
	return &Manager{mc}
}

func NewManagerDefault() *Manager {
	return NewManager(NewConfigDefault())
}

func (m *Manager) Recording() (*Logger, func()) {
	l := NewLogger(m.Config)
	flush := l.Start()
	return l, flush
}

func (m *Manager) RecordingWithContext(ctx context.Context) (context.Context, func()) {
	l, flush := m.Recording()
	nctx := NewContextWithLogger(ctx, l)
	return nctx, flush
}
