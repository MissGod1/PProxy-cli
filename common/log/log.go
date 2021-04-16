package log

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	NONE
)

var level = INFO
var logger = NewLogger()

func SetLevel(_level int) {
	level = _level
}

func Debugf(msg string, args ...interface{}) {
	if level <= DEBUG {
		logger.Printf("[ DEBUG ] "+msg, args)
	}
}

func Infof(msg string, args ...interface{}) {
	if level <= INFO {
		logger.Printf("[ INFO  ] "+msg, args)
	}
}

func Errorf(msg string, args ...interface{}) {
	if level <= ERROR {
		logger.Printf("[ ERROR ] "+msg, args)
	}
}

func Warningf(msg string, args ...interface{}) {
	if level <= WARNING {
		logger.Printf("[WARNING] "+msg, args)
	}
}

func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args)
}
