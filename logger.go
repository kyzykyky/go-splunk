package gosplunk

import "fmt"

type Logger interface {
	Info(args ...interface{})
	Infow(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnw(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorw(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalw(args ...interface{})
	Fatalf(template string, args ...interface{})
	Debug(args ...interface{})
	Debugw(args ...interface{})
	Debugf(template string, args ...interface{})
}

type NoLogger struct{}

func (n NoLogger) Info(args ...interface{})                    {}
func (n NoLogger) Infow(args ...interface{})                   {}
func (n NoLogger) Infof(template string, args ...interface{})  {}
func (n NoLogger) Warn(args ...interface{})                    {}
func (n NoLogger) Warnw(args ...interface{})                   {}
func (n NoLogger) Warnf(template string, args ...interface{})  {}
func (n NoLogger) Error(args ...interface{})                   {}
func (n NoLogger) Errorw(args ...interface{})                  {}
func (n NoLogger) Errorf(template string, args ...interface{}) {}
func (n NoLogger) Fatal(args ...interface{})                   {}
func (n NoLogger) Fatalw(args ...interface{})                  {}
func (n NoLogger) Fatalf(template string, args ...interface{}) {}
func (n NoLogger) Debug(args ...interface{})                   {}
func (n NoLogger) Debugw(args ...interface{})                  {}
func (n NoLogger) Debugf(template string, args ...interface{}) {}

// SimpleLogger is a simple logger that prints to stdout
type SimpleLogger struct{}

func (s SimpleLogger) Info(args ...interface{}) {
	fmt.Print("INFO:")
	fmt.Println(args...)
}
func (s SimpleLogger) Infow(args ...interface{}) {
	fmt.Print("INFO:")
	fmt.Println(args...)
}
func (s SimpleLogger) Infof(template string, args ...interface{}) {
	fmt.Print("INFO:")
	fmt.Printf(template, args...)
}
func (s SimpleLogger) Warn(args ...interface{}) {
	fmt.Print("WARN:")
	fmt.Println(args...)
}
func (s SimpleLogger) Warnw(args ...interface{}) {
	fmt.Print("WARN:")
	fmt.Println(args...)
}
func (s SimpleLogger) Warnf(template string, args ...interface{}) {
	fmt.Print("WARN:")
	fmt.Printf(template, args...)
}
func (s SimpleLogger) Error(args ...interface{}) {
	fmt.Print("ERROR:")
	fmt.Println(args...)
}
func (s SimpleLogger) Errorw(args ...interface{}) {
	fmt.Print("ERROR:")
	fmt.Println(args...)
}
func (s SimpleLogger) Errorf(template string, args ...interface{}) {
	fmt.Print("ERROR:")
	fmt.Printf(template, args...)
}
func (s SimpleLogger) Fatal(args ...interface{}) {
	fmt.Print("FATAL:")
	fmt.Println(args...)
	panic(args)
}
func (s SimpleLogger) Fatalw(args ...interface{}) {
	fmt.Print("FATAL:")
	fmt.Println(args...)
	panic(args)
}
func (s SimpleLogger) Fatalf(template string, args ...interface{}) {
	fmt.Print("FATAL:")
	fmt.Printf(template, args...)
	panic(args)
}
func (s SimpleLogger) Debug(args ...interface{}) {
	fmt.Print("DEBUG:")
	fmt.Println(args...)
}
func (s SimpleLogger) Debugw(args ...interface{}) {
	fmt.Print("DEBUG:")
	fmt.Println(args...)
}
func (s SimpleLogger) Debugf(template string, args ...interface{}) {
	fmt.Print("DEBUG:")
	fmt.Printf(template, args...)
}
