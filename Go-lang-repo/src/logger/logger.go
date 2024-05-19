package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logger *log.Logger

//Use the below Multiliine code to send logs to "logs/server.log" file. For now the logs are being shared to docker console

/*

func InitLogger() *log.Logger {
	logFilePath := os.Getenv("LOG_FILE_PATH") 
	if logFilePath == "" {
		logFilePath = "./logs/server.log" 
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}


	// Load the IST location
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatal("Error loading IST location:", err)
	}

	// Create a new logger with IST timezone
	logger = log.New(logFile, "", 0)
	logger.SetFlags(log.Lshortfile)
	logger.SetPrefix(time.Now().In(loc).Format("02/01/2006 15:04:05") + " ")

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

*/

func InitLogger() *log.Logger {
	logFilePath := os.Getenv("LOG_FILE_PATH") 
	if logFilePath == "" {
		logFilePath = "./logs/server.log" 
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}

	// Load the IST location
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatal("Error loading IST location:", err)
	}

	// Create a new logger with IST timezone
	logger = log.New(logFile, "", 0)
	logger.SetFlags(log.Lshortfile)
	logger.SetPrefix(time.Now().In(loc).Format("02/01/2006 15:04:05") + " ")

	return logger
}

func Log(v ...interface{}) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatal("Error loading IST location:", err)
	}

	if logger != nil {
		timestamp := time.Now().In(loc).Format("02/01/2006 15:04:05")
		fmt.Print(timestamp + " ")
		fmt.Println(v...)
	}
}

func Printf(format string, v ...interface{}) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatal("Error loading IST location:", err)
	}

	if logger != nil {
		timestamp := time.Now().In(loc).Format("02/01/2006 15:04:05")
		fmt.Printf(timestamp + " " + format, v...)
	}
}
