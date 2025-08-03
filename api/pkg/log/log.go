package log

type options struct {
	logLevel          string
	sentryDSN         string
	sentryServerName  string
	sentryEnvironment string
	sentryLevel       string
}

type Option func(opts *options)

func buildOptions(opts ...Option) *options {
	dopts := &options{
		logLevel:          "info",
		sentryDSN:         "",
		sentryServerName:  "",
		sentryEnvironment: "",
		sentryLevel:       "warn",
	}
	for i := range opts {
		opts[i](dopts)
	}
	return dopts
}

func WithLogLevel(level string) Option {
	return func(opts *options) {
		opts.logLevel = level
	}
}

func WithSentryDSN(dsn string) Option {
	return func(opts *options) {
		opts.sentryDSN = dsn
	}
}

func WithSentryServerName(name string) Option {
	return func(opts *options) {
		opts.sentryServerName = name
	}
}

func WithSentryEnvironment(env string) Option {
	return func(opts *options) {
		opts.sentryEnvironment = env
	}
}

func WithSentryLevel(level string) Option {
	return func(opts *options) {
		opts.sentryLevel = level
	}
}
