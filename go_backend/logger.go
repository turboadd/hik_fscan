package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	logFile       *os.File
	logConsole    io.Writer = os.Stdout
	logFileWriter io.Writer
	currentLevel  = LevelDebug
)

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
)

func InitLogger(logPath string, level int) error {

	currentLevel = level

	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	logFileWriter = &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    5,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	}
	return nil
}

func logWithColor(levelStr, color, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	consoleMsg := fmt.Sprintf("[%s] %s[%s]%s %s", timestamp, color, levelStr, colorReset, msg)
	fmt.Fprintln(logConsole, consoleMsg)

	fileMsg := fmt.Sprintf("[%s] [%s] %s\n", timestamp, levelStr, msg)
	fmt.Fprint(logFileWriter, fileMsg)
}

func Debug(msg string) {
	if currentLevel <= LevelDebug {
		logWithColor("DEBUG", colorBlue, msg)
	}
}

func Info(msg string) {
	if currentLevel <= LevelInfo {
		logWithColor("INFO", colorGreen, msg)
	}
}

func Warn(msg string) {
	if currentLevel <= LevelWarn {
		logWithColor("WARN", colorYellow, msg)
	}
}

func Error(msg string) {
	if currentLevel <= LevelError {
		logWithColor("ERROR", colorRed, msg)
	}
}
