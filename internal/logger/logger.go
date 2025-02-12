package logger

type Logger interface {
	Debugf(template string, args ...any)
	Infof(template string, args ...any)
	Errorf(template string, args ...any)
	Fatalf(template string, args ...any)
}
