package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "", 0)
}

func LogInfo(format string, v ...interface{}) {
	logWithLevel("INFO", format, v...)
}

func LogError(format string, v ...interface{}) {
	logWithLevel("ERROR", format, v...)
}

func LogWarning(format string, v ...interface{}) {
	logWithLevel("WARN", format, v...)
}

func logWithLevel(level, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)
	Logger.Printf("[%s] %s | %s", level, timestamp, message)
}
