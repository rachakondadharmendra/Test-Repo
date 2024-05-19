package logger

import (
	"log"
	"os"
)

var logger *log.Logger

func InitLogger() *log.Logger {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

func Log(v ...interface{}) {
	if logger != nil {
		logger.Println(v...)
	}
}

func Printf(format string, v ...interface{}) {
	if logger != nil {
		logger.Printf(format, v...)
	}
}
