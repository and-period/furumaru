package log

import (
	"context"
	"errors"
	"log/slog"
	"runtime"
	"strconv"

	"github.com/getsentry/sentry-go"
)

const (
	sentryLogger      = "slog"
	sentryContextKey  = "attributes"
	sentryMaxErrDepth = 10
	levelFatal        = slog.Level(12)
)

var _ slog.Handler = (*sentryHandler)(nil)

// sentryHandler は指定レベル以上のログを Sentry の Issue（イベント）として送信する slog.Handler。
// sentry-go v0.48 で slog 連携からのイベント送信が削除されたため、自前で実装している。
type sentryHandler struct {
	hub      *sentry.Hub
	minLevel slog.Level
	attrs    []slog.Attr
	groups   []string
}

func newSentryHandler(ctx context.Context, minLevel slog.Level) *sentryHandler {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}

	return &sentryHandler{hub: hub, minLevel: minLevel}
}

func (h *sentryHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.minLevel
}

func (h *sentryHandler) Handle(ctx context.Context, record slog.Record) error {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = h.hub
	}

	attrs := make(map[string]any, record.NumAttrs()+len(h.attrs))

	var errs []error
	// WithAttrs で保持した属性はグループ適用済み
	for _, a := range h.attrs {
		h.collectAttr(attrs, &errs, nil, a)
	}

	record.Attrs(func(a slog.Attr) bool {
		h.collectAttr(attrs, &errs, h.groups, a)
		return true
	})

	event := sentry.NewEvent()
	event.Level = sentryLevel(record.Level)
	event.Message = record.Message

	event.Logger = sentryLogger
	if !record.Time.IsZero() {
		event.Timestamp = record.Time
	}

	if record.PC != 0 {
		frame, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
		attrs["source"] = frame.File + ":" + strconv.Itoa(frame.Line)
	}

	if len(attrs) > 0 {
		event.Contexts[sentryContextKey] = attrs
	}

	switch len(errs) {
	case 0:
	case 1:
		event.SetException(errs[0], sentryMaxErrDepth)
	default:
		event.SetException(errors.Join(errs...), sentryMaxErrDepth)
	}

	hub.CaptureEvent(event)

	return nil
}

func (h *sentryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	nh := *h
	nh.attrs = make([]slog.Attr, 0, len(h.attrs)+len(attrs))

	nh.attrs = append(nh.attrs, h.attrs...)
	for _, a := range attrs {
		for i := len(h.groups) - 1; i >= 0; i-- {
			a = slog.Group(h.groups[i], a)
		}

		nh.attrs = append(nh.attrs, a)
	}

	return &nh
}

func (h *sentryHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	nh := *h
	nh.groups = append(append([]string{}, h.groups...), name)

	return &nh
}

func (h *sentryHandler) collectAttr(
	dst map[string]any,
	errs *[]error,
	groups []string,
	a slog.Attr,
) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}

	key := a.Key
	for i := len(groups) - 1; i >= 0; i-- {
		key = groups[i] + "." + key
	}

	if a.Value.Kind() == slog.KindGroup {
		sub := append(append([]string{}, groups...), a.Key)
		if a.Key == "" {
			sub = groups
		}

		for _, ga := range a.Value.Group() {
			h.collectAttr(dst, errs, sub, ga)
		}

		return
	}

	if err, ok := a.Value.Any().(error); ok {
		*errs = append(*errs, err)
		dst[key] = err.Error()

		return
	}

	dst[key] = a.Value.Any()
}

func sentryLevel(level slog.Level) sentry.Level {
	switch {
	case level >= levelFatal:
		return sentry.LevelFatal
	case level >= slog.LevelError:
		return sentry.LevelError
	case level >= slog.LevelWarn:
		return sentry.LevelWarning
	case level >= slog.LevelInfo:
		return sentry.LevelInfo
	default:
		return sentry.LevelDebug
	}
}
