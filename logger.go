package logger

type formatter interface {
	format(string) string
}

type prefixFormatter struct {
	prefix string
}

func (pf prefixFormatter) format(s string) string {
	return pf.prefix + s
}

type suffixFormatter struct {
	suffix string
}

func (sf suffixFormatter) format(s string) string {
	return s + sf.suffix
}

type Logger interface {
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}

type fmtLogger struct {
	l Logger
	f formatter
}

func (fl fmtLogger) Infof(format string, a ...interface{}) {
	fl.l.Infof(fl.f.format(format), a...)
}

func (fl fmtLogger) Warnf(format string, a ...interface{}) {
	fl.l.Warnf(fl.f.format(format), a...)
}

func (fl fmtLogger) Errorf(format string, a ...interface{}) {
	fl.l.Errorf(fl.f.format(format), a...)
}

func WithPrefix(logger Logger, prefix string) Logger {
	return fmtLogger{l: logger, f: prefixFormatter{prefix: prefix}}
}

func WithSuffix(logger Logger, suffix string) Logger {
	return fmtLogger{l: logger, f: suffixFormatter{suffix: suffix}}
}
