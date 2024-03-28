package logger

// Info func
func Info(args ...interface{}) {
	global.Info(args)
}

// Infof func
func Infof(tmp string, args ...interface{}) {
	global.Infof(tmp, args...)
}

// Fatal func
func Fatal(args ...interface{}) {
	global.Fatal(args)
}

// Fatalf func
func Fatalf(tmp string, args ...interface{}) {
	global.Fatalf(tmp, args...)
}

// Error func
func Error(args ...interface{}) {
	global.Error(args)
}

// Errorf func
func Errorf(tmp string, args ...interface{}) {
	global.Errorf(tmp, args...)
}
