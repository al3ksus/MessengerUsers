package logger

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=Logger
type Logger interface {
	Debugf(template string, args ...any)
	Infof(template string, args ...any)
	Warnf(template string, args ...any)
	Errorf(template string, args ...any)
}
