package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	cfg "github.com/satont/twir/libs/config"
)

type Logger interface {
	Info(input string, fields ...any)
	Error(input string, fields ...any)
	Debug(input string, fields ...any)
	WithComponent(name string) Logger
}

type logger struct {
	log *slog.Logger

	service string
	sentry  *sentry.Client
}

type Opts struct {
	Env     string
	Service string

	Sentry *sentry.Client
}

func NewFx(opts Opts) func(config cfg.Config, sentry *sentry.Client) Logger {
	return func(config cfg.Config, sentry *sentry.Client) Logger {
		return New(
			Opts{
				Env:     config.AppEnv,
				Service: opts.Service,
				Sentry:  sentry,
			},
		)
	}
}

func New(opts Opts) Logger {
	var log *slog.Logger

	switch opts.Env {
	case "development":
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug, AddSource: true,
				},
			),
		)
	case "production":
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true},
			),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true},
			),
		)
	}

	if opts.Service != "" {
		log = log.With(slog.String("service", opts.Service))
	}

	return &logger{
		log:     log,
		sentry:  opts.Sentry,
		service: opts.Service,
	}
}

func (c *logger) handle(level slog.Level, input string, fields ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, input, pcs[0])
	for _, f := range fields {
		r.Add(f)
	}
	_ = c.log.Handler().Handle(context.Background(), r)
}

func (c *logger) Info(input string, fields ...any) {
	c.handle(slog.LevelInfo, input, fields...)
}

func (c *logger) Error(input string, fields ...any) {
	c.handle(slog.LevelError, input, fields...)

	if c.sentry != nil {
		scope := sentry.NewScope()

		for _, f := range fields {
			casted, ok := f.(slog.Attr)
			if !ok {
				continue
			}

			scope.SetExtra(casted.Key, casted.Value.Any())
		}

		scope.SetTag("service", c.service)
		scope.SetExtra("service", c.service)

		c.sentry.CaptureMessage(input, &sentry.EventHint{}, scope)
	}
}

func (c *logger) Debug(input string, fields ...any) {
	c.handle(slog.LevelDebug, input, fields...)
}

func (c *logger) WithComponent(name string) Logger {
	return &logger{
		log:     c.log.With(slog.String("component", name)),
		sentry:  c.sentry,
		service: c.service,
	}
}
