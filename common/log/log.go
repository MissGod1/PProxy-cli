package log

import "log"

const (
	DEBUG   = iota
	INFO
	ERROR
	WARNING
	NONE
)

var level = INFO

func SetLevel(_level int) {
	level = _level
}

func Debugf(msg string, args ...interface{}) {
	if level <= DEBUG {
		log.Printf("[ DEBUG ] " + msg, args)
	}
}

func Infof(msg string, args ...interface{}) {
	if level <= INFO {
		log.Printf("[ INFO  ] " + msg, args)
	}
}

func Errorf(msg string, args ...interface{}) {
	if level <= ERROR {
		log.Printf("[ ERROR ] " + msg, args)
	}
}

func Warningf(msg string, args ...interface{}) {
	if level <= WARNING {
		log.Printf("[WARNING] " + msg, args)
	}
}

func Fatalf(msg string, args ...interface{}) {
	log.Fatalf(msg, args)
}