package console

import (
	"fmt"
	"io"
)

const (
	moveCursorOneLineUpInTheTerminal = "\033[A"
	moveCursorToStartOfLine          = "\r"
	moveCursorToEndOfLine            = "\033[999C"
	clearConsole                     = "\033[H\033[2J"
	clearLine                        = "\033[2K"
	newLine                          = "\r\n"
)

var lastPrintWasASpinner = false
var spinnerCount int = 0 //Overflow is fine
// var indicators = []string{"-", "\\", "|", "/"}
var indicators = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

type commonArguments struct {
	Types bool
}

type GoIsDumb struct{}

var Types = GoIsDumb{}

func PrettyMemoryAllocation(memoryInBytes uint64) string {
	if memoryInBytes >= 1<<20 { // 1 MB or more
		return fmt.Sprintf("%.2f MB", float64(memoryInBytes)/(1<<20))
	} else if memoryInBytes >= 1<<10 { // 1 KB or more
		return fmt.Sprintf("%.2f KB", float64(memoryInBytes)/(1<<10))
	} else {
		return fmt.Sprintf("%d bytes", memoryInBytes)
	}
}

func ShouldLog(logLevel LogLevel, minimumLogLevel LogLevel) bool {
	if minimumLogLevel == NOT_SPECIFIED {
		return true
	}
	return logLevel >= minimumLogLevel
}

func MaximumLogLevel(logLevel LogLevel, maximumLogLevel LogLevel) bool {
	if maximumLogLevel == NOT_SPECIFIED {
		return true
	}
	return logLevel < maximumLogLevel
}

func (GoIsDumb) Log(arguments ...any) {
	lastPrintWasASpinner = false
	defaultLogger.common(INFO, commonArguments{Types: true}, arguments...)
}

func Spinner(arguments ...any) {
	spinnerCount++
	controlCharacters := ""
	if lastPrintWasASpinner {
		controlCharacters = moveCursorOneLineUpInTheTerminal + moveCursorToStartOfLine + clearLine
	}
	nextIndicator := indicators[spinnerCount%len(indicators)]

	defaultLogger.Write(DEBUG, []byte(controlCharacters+nextIndicator+" "))
	defaultLogger.common(DEBUG, commonArguments{Types: false}, arguments...)
	defaultLogger.Write(DEBUG, []byte(moveCursorOneLineUpInTheTerminal+moveCursorToEndOfLine+nextIndicator+newLine))
	lastPrintWasASpinner = true
}

func Clear() {
	defaultLogger.Clear()
}

func Copy(src io.Reader) {
	defaultLogger.Copy(src)
}

func Log(arguments ...any) {
	lastPrintWasASpinner = false
	defaultLogger.Log(arguments...)
}

func Info(arguments ...any) {
	lastPrintWasASpinner = false
	defaultLogger.Info(arguments)
}

func Debug(arguments ...any) {
	lastPrintWasASpinner = false
	defaultLogger.Debug(arguments)
}

func Trace(arguments ...any) {
	lastPrintWasASpinner = false
	defaultLogger.Trace(arguments)
}

func Error(arguments ...any) {
	lastPrintWasASpinner = false
	defaultLogger.Error(arguments)
}

func Fatal(arguments ...any) {
	lastPrintWasASpinner = false
	defaultLogger.Fatal(arguments...)
}

func (log *Logger) Clear() {
	log.Write(INFO, []byte(clearConsole))
}

func (log *Logger) Copy(src io.Reader) {
	log.LogCopy(DEBUG, src)
}

func (log *Logger) Log(arguments ...any) {
	log.Info(arguments...)
}

func (log *Logger) Info(arguments ...any) {
	log.common(INFO, commonArguments{Types: false}, arguments...)
}

func (log *Logger) Debug(arguments ...any) {
	log.common(DEBUG, commonArguments{Types: false}, arguments...)
}

func (log *Logger) Trace(arguments ...any) {
	log.common(TRACE, commonArguments{Types: false}, arguments...)
}

func (log *Logger) Error(arguments ...any) {
	log.common(ERROR, commonArguments{Types: false}, arguments...)
}

func (log *Logger) Fatal(arguments ...any) {
	log.common(ERROR, commonArguments{Types: false}, arguments...)
	panic("Fatal")
}
