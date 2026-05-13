// Package logger provides the application-wide logging interface backed by golang-wappers.
package logger

import (
	wappers "github.com/FrancoPersonal/golang-wappers/pkg/logger"
)

// Logger is the application-wide logging interface.
// Domain and application packages depend only on this interface, never on the concrete implementation.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	With(args ...any) Logger
}

// wappersAdapter adapts github.com/FrancoPersonal/golang-wappers/pkg/logger to the local Logger interface.
type wappersAdapter struct {
	inner wappers.Logger
}

// New returns a Logger backed by golang-wappers.
// level accepts: "debug", "info", "warn", "error".
func New(level string) Logger {
	var lvl wappers.Level
	switch level {
	case "debug":
		lvl = wappers.LevelDebug
	case "warn":
		lvl = wappers.LevelWarn
	case "error":
		lvl = wappers.LevelError
	default:
		lvl = wappers.LevelInfo
	}
	return &wappersAdapter{
		inner: wappers.New(wappers.Options{Level: lvl, JSON: true}),
	}
}

func (a *wappersAdapter) Info(msg string, args ...any)  { a.inner.Info(msg, args...) }
func (a *wappersAdapter) Error(msg string, args ...any) { a.inner.Error(msg, args...) }
func (a *wappersAdapter) Debug(msg string, args ...any) { a.inner.Debug(msg, args...) }
func (a *wappersAdapter) Warn(msg string, args ...any)  { a.inner.Warn(msg, args...) }

func (a *wappersAdapter) With(args ...any) Logger {
	return &wappersAdapter{inner: a.inner.With(args...)}
}
