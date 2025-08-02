package log

import (
	"log/slog"
)

func Error(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Any("error", err)
}

func Strings(key string, values []string) slog.Attr {
	return slog.Any(key, values)
}

func Ints(key string, values []int) slog.Attr {
	return slog.Any(key, values)
}

func Int32s(key string, values []int32) slog.Attr {
	return slog.Any(key, values)
}

func Int64s(key string, values []int64) slog.Attr {
	return slog.Any(key, values)
}
