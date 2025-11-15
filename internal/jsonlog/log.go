package jsonlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Level Define a Level type to represent the severity level for a log entry.
type Level int8

// Logger Define a custom Logger type. This holds the output destination that the log entries
// will be written to, the minimum severity level that log entries will be written for,
// plus a mutex for coordinating the writes.
type Logger struct {
	out      io.Writer
	minLevel Level
	mutex    sync.Mutex
}

const (
	LevelInfo Level = iota // Has the value 0
	LevelInfoError
	LevelError
	LevelFatal
	LevelOff
)

// Return a human-friendly string for the severity level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelInfoError:
		return "ERROR"
	case LevelError:
		return "ERROR+STACK"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

var (
	log *Logger
)

func init() {
	// Initialize a new jsonlog.Logger which writes any messages *at or above* the INFO severity level to the standard out stream.
	log = New(os.Stdout, LevelInfo)
}

// New Return a new Logger instance which writes log entries at or above a minimum severity
// level to a specific output destination.
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// MARK: - Info
func Info(format string, args ...interface{}) {
	var message string
	if len(args) > 0 {
		message = fmt.Sprintf("💭 "+format, args...)
	} else {
		message = "💭 " + format
	}
	writeLog(LevelInfo, message, nil)
}

// Info Declare some helper methods for writing log entries at the different levels. Notice
// that these all accept a map as the second parameter which can contain any arbitrary
// 'properties' that you want to appear in the log entry.
func InfoWithProperties(message string, properties map[string]string) {
	writeLog(LevelInfo, "💭 "+message, properties)
}

// MARK: - Error
func Error(format string, args ...interface{}) {
	message := fmt.Sprintf("❌ "+format, args...)
	writeLog(LevelInfoError, message, nil)
}

func ErrorWithProperties(err error, properties map[string]string) {
	writeLog(LevelError, "❌ "+err.Error(), properties)
}

// MARK: - Fatal
func Fatal(err error) {
	writeLog(LevelFatal, "🆘 "+err.Error(), nil)
	os.Exit(1) // For entries at the FATAL level, we also terminate the application.
}

func FatalWithProperties(err error, properties map[string]string) {
	writeLog(LevelFatal, "🆘 "+err.Error(), properties)
	os.Exit(1) // For entries at the FATAL level, we also terminate the application.
}

func writeLog(level Level, message string, properties map[string]string) (int, error) {
	// If the severity level of the log entry is below the minimum severity for the
	// logger, then return with no further action.
	if level < log.minLevel {
		return 0, nil
	}

	// Declare an anonymous struct holding the data for the log entry.
	aux := struct {
		Level      string            `json:"level"`
		Time       string            `json:"time"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
		Trace      string            `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().In(time.FixedZone("PST", -8*60*60)).Format("02-Jan-06 15:04:05.999 MST"),
		Message:    message,
		Properties: properties,
	}

	// Include a stack trace for entries at the ERROR and FATAL levels.
	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	// Declare a line variable for holding the actual log entry text.
	var line []byte

	// Marshal the anonymous struct to JSON and store it in the line variable. If there
	// was a problem creating the JSON, set the contents of the log entry to be that
	// plain-text error message instead.
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	}

	// Lock the mutex so that no two writes to the output destination can happen
	// concurrently. If we don't do this, it's possible that the text for two or more
	// log entries will be intermingled in the output.
	log.mutex.Lock()
	defer log.mutex.Unlock()

	// Write the log entry followed by a newline.
	return log.out.Write(append(line, '\n'))
}

// We also implement a Write() method on our Logger type so that it satisfies the
// io.Writer interface. This writes a log entry at the ERROR level with no additional
// properties.
func (l *Logger) Write(message []byte) (n int, err error) {
	return writeLog(LevelError, string(message), nil)
}
