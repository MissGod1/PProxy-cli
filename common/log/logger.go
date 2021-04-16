package log

import (
	"io"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	logFile, err := os.OpenFile("pproxy-cli.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败:", err)
	}

	return log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)
}
