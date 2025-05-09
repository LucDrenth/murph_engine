package log

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/lucdrenth/murph_engine/src/log/ansi"
)

type ConsoleLogger struct {
	DebugColor ansi.Color
	InfoColor  ansi.Color
	WarnColor  ansi.Color
	ErrorColor ansi.Color
	TraceColor ansi.Color

	TimestampColor  ansi.Color
	TimestampFormat string

	CallerColor     ansi.Color
	LogCaller       bool  // if true, log line include the file and line number of the method that called the log method
	CallerPathDepth int   // the number of caller directories to include. For example, use 3 for path/to/file.go
	Level           Level // skip logs that are lower than this level

	storage Storage
}

var _ Logger = &ConsoleLogger{}

// Console returns a logger that prints colored output to the console. Output format looks like this:
//
// 12:30:01 | INFO | main.go:15 - a message with a nice color
func Console() ConsoleLogger {
	return ConsoleLogger{
		DebugColor:      ansi.ColorGrey,
		InfoColor:       ansi.ColorWhite,
		WarnColor:       ansi.ColorYellow,
		ErrorColor:      ansi.ColorBrightRed,
		TraceColor:      ansi.ColorGrey,
		CallerColor:     ansi.ColorCyan,
		TimestampColor:  ansi.ColorGreen,
		TimestampFormat: "15:04:05.0000",
		LogCaller:       true,
		CallerPathDepth: 3,
		Level:           LevelDebug,
		storage:         NewStorage(),
	}
}

func (logger *ConsoleLogger) Debug(message string) {
	logger.log(LevelDebug, logger.DebugColor, message, 3)
}

func (logger *ConsoleLogger) DebugOnce(message string) {
	logger.logOnce(LevelDebug, logger.DebugColor, message)
}

func (logger *ConsoleLogger) Info(message string) {
	logger.log(LevelInfo, logger.InfoColor, message, 3)
}

func (logger *ConsoleLogger) InfoOnce(message string) {
	logger.logOnce(LevelInfo, logger.InfoColor, message)
}

func (logger *ConsoleLogger) Warn(message string) {
	logger.log(LevelWarn, logger.WarnColor, message, 3)
}

func (logger *ConsoleLogger) WarnOnce(message string) {
	logger.logOnce(LevelWarn, logger.WarnColor, message)
}

func (logger *ConsoleLogger) Error(message string) {
	logger.log(LevelError, logger.ErrorColor, message, 3)
}

func (logger *ConsoleLogger) ErrorOnce(message string) {
	logger.logOnce(LevelError, logger.ErrorColor, message)
}

func (logger *ConsoleLogger) Trace(message string) {
	stackTrace := string(debug.Stack())
	message = fmt.Sprintf("%s\n%s", message, stackTrace)

	logger.log(levelStackTrace, logger.TraceColor, message, 3)
}

func (logger *ConsoleLogger) TraceOnce(message string) {
	stackTrace := string(debug.Stack())
	message = fmt.Sprintf("%s\n%s", message, stackTrace)

	logger.logOnce(levelStackTrace, logger.TraceColor, message)
}

func (logger *ConsoleLogger) logOnce(level Level, messageColor ansi.Color, message string) {
	caller, ok := logger.getCallerForStorage()
	if !ok {
		// We unexpectedly failed to get the caller.
		// Let's try to continue to log our original message with an empty caller and log an additional error message.
		logger.Error("failed to get caller location")
	}

	if logger.storage.Exists(level, message, caller) {
		return
	}

	didLog := logger.log(level, messageColor, message, 4)
	if !didLog {
		return
	}

	logger.storage.Insert(level, message, caller)
}

// log returns wether the message was logged or not.
// The message does not get logged if the log level does not allow for it.
func (logger *ConsoleLogger) log(level Level, messageColor ansi.Color, message string, callerDepth int) bool {
	if !logger.Level.Allows(level) {
		return false
	}

	caller, ok := logger.getCallerForLogMessage(callerDepth)
	if ok {
		caller = fmt.Sprintf("| %s", ansi.Colorize(logger.CallerColor, caller))
	}

	fmt.Printf("%s | %-14s %s - %s\n",
		ansi.Colorize(logger.TimestampColor, time.Now().Format(logger.TimestampFormat)),
		ansi.Colorize(messageColor, strings.ToUpper(level.String())),
		caller,
		ansi.Colorize(messageColor, message),
	)

	return true
}

func (logger *ConsoleLogger) getCallerForStorage() (caller string, ok bool) {
	_, fullPath, line, ok := runtime.Caller(3) // 3 means we get the caller of the log methods (Info, InfoOnce etc.)
	return fmt.Sprintf("%s:%d", fullPath, line), ok
}

func (logger *ConsoleLogger) getCallerForLogMessage(callerDepth int) (caller string, ok bool) {
	if !logger.LogCaller {
		return "", false
	}

	_, fullPath, line, ok := runtime.Caller(callerDepth)
	if !ok {
		return "", false
	}

	var path string

	pathSplit := strings.Split(fullPath, "/")
	if len(pathSplit) <= logger.CallerPathDepth {
		path = strings.Join(pathSplit, "/")
	} else {
		path = strings.Join(pathSplit[len(pathSplit)-logger.CallerPathDepth:], "/")
	}

	return fmt.Sprintf("%s:%d", path, line), true
}

func (logger *ConsoleLogger) ClearStorage() {
	logger.storage.Clear()
}
