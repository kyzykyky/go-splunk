package logger

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnw(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorw(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
}

type NoLogger struct{}

func (n NoLogger) Info(args ...interface{}) {
}
func (n NoLogger) Infof(template string, args ...interface{}) {
}
func (n NoLogger) Warn(args ...interface{}) {
}
func (n NoLogger) Warnw(template string, args ...interface{}) {
}
func (n NoLogger) Warnf(template string, args ...interface{}) {
}
func (n NoLogger) Error(args ...interface{}) {
}
func (n NoLogger) Errorw(template string, args ...interface{}) {
}
func (n NoLogger) Errorf(template string, args ...interface{}) {
}
func (n NoLogger) Fatal(args ...interface{}) {
}
func (n NoLogger) Fatalf(template string, args ...interface{}) {
}
func (n NoLogger) Debug(args ...interface{}) {
}
func (n NoLogger) Debugf(template string, args ...interface{}) {
}
