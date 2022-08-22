package log

import (
	"context"
	"errors"
	"io"

	"github.com/rs/zerolog/hlog"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info(ctx context.Context, msg string, params ...param)
	Debug(ctx context.Context, msg string, params ...param)
	Warn(ctx context.Context, msg string, params ...param)
	Err(ctx context.Context, msg string, err error, params ...param)
	Fatal(ctx context.Context, msg string, err error, params ...param)
}

type logger struct {
	Log zerolog.Logger
}

func New(writer io.Writer, isJSON bool) logger {
	const timeFormat = "02 Jan 2006 15:04:05"

	if !isJSON {
		writer = zerolog.ConsoleWriter{Out: writer, TimeFormat: timeFormat}
	}

	return logger{
		Log: zerolog.New(writer).With().Timestamp().Logger(),
	}
}

type param struct {
	Key   string
	Value any
}

func P(key string, value any) param {
	return param{key, value}
}

func leveledEvent(ctx context.Context, evt *zerolog.Event, msg string, params ...param) *zerolog.Event {
	evt = evt.Caller(2)

	for _, p := range params {
		evt.Interface(p.Key, p.Value)
	}

	if logID, ok := hlog.IDFromCtx(ctx); ok {
		evt.Stringer("request_id", logID)
	}

	return evt
}

func (l logger) Info(ctx context.Context, msg string, params ...param) {
	leveledEvent(ctx, l.Log.Info(), msg, params...).Msg(msg)
}

func (l logger) Debug(ctx context.Context, msg string, params ...param) {
	leveledEvent(ctx, l.Log.Debug(), msg, params...).Msg(msg)
}

func (l logger) Warn(ctx context.Context, msg string, params ...param) {
	leveledEvent(ctx, l.Log.Warn(), msg, params...).Msg(msg)
}

func (l logger) Err(ctx context.Context, msg string, err error, params ...param) {
	if err == nil {
		err = errors.New("nil")
	}

	leveledEvent(ctx, l.Log.Err(err), msg, params...).Msg(msg)
}

func (l logger) Fatal(ctx context.Context, msg string, err error, params ...param) {
	leveledEvent(ctx, l.Log.Fatal(), msg, params...).Err(err).Msg(msg)
}
