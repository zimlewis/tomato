package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"slices"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	pkgerr "github.com/pkg/errors"
	"github.com/zimlewis/tomato/internal/tomatoerrs"
)

type closeFunc func() error
type multiError interface {
	error
	Unwrap() []error
}
type stackTracer interface {
	error
	StackTrace() pkgerr.StackTrace
}

// Initialize a logger that is configured to log structured error
func initializeLogger() (*slog.Logger, closeFunc, error) {
	cf := func () error { return nil }
	var handlers []slog.Handler
	debugHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: replaceAttr,
	})
	handlers = append(handlers, debugHandler)

	multiHandler := slog.NewMultiHandler(handlers...)
	logger := slog.New(multiHandler)
	return logger, cf, nil
}


func errorToArgs(err error) []slog.Attr {
	attrs := tomatoerrs.Attrs(err)
	attrs = append(attrs, slog.Attr{
		Key: "message",
		Value: slog.StringValue(err.Error()),
	})

	if stackErr, ok := errors.AsType[stackTracer](err); ok {
		attrs = append(
			attrs,
			slog.Attr{
				Key:   "stack_trace",
				Value: slog.StringValue(fmt.Sprintf("%+v", stackErr.StackTrace())),
			},
		)
	}

	return attrs
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	var sensitiveKeys = []string{"password", "key", "apikey", "secret", "pin", "creditcardno", "user"}
	if a.Key == "error" {

		err, ok := a.Value.Any().(error)
		if !ok {
			return a
		}

		if errs, ok := errors.AsType[multiError](err); ok {
			var attrs []slog.Attr
			for i, err := range errs.Unwrap() {
				attrs = append(
					attrs, 
					slog.GroupAttrs(
						fmt.Sprintf("error_%d", i + 1), 
						errorToArgs(err)...,
					),
				)
			}
			return slog.GroupAttrs("errors", attrs...)
		} else {
			var attrs []slog.Attr
			attrs = append(attrs, errorToArgs(err)...)
			return slog.GroupAttrs("error", attrs...)
		}

	}
	if slices.Contains(sensitiveKeys, a.Key) {
		return slog.String(a.Key, "[REDACTED]")
	}

	if u, ok := a.Value.Any().(string); ok {
		parsed, err := url.Parse(u)
		if err != nil {	
			goto skip
		}
		_, ok := parsed.User.Password()
		if !ok {
			goto skip
		}

		parsed.User = url.UserPassword(parsed.User.Username(), "[REDACTED]")
		return slog.String(a.Key, parsed.String())
	}
	skip:

	return a
}

func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
