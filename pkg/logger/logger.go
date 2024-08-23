package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func LogHttp(h func(w http.ResponseWriter, r *http.Request, _ httprouter.Params)) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var start time.Time = time.Now()
		h(w, r, ps)
		fmt.Printf("\nMETHOD: %s, TIME TAKEN: %d\n", r.Method, time.Since(start).Milliseconds())
	}
}

package advancedlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[LogLevel]string{
	TRACE: "TRACE",
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Fields    []Field
	Caller    string
}

// Formatter interface for custom log formatters
type Formatter interface {
	Format(entry *LogEntry) ([]byte, error)
}

// TextFormatter formats log entries as plain text
type TextFormatter struct{}

func (f *TextFormatter) Format(entry *LogEntry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] %s %s: %s %v\n",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		levelNames[entry.Level],
		entry.Caller,
		entry.Message,
		entry.Fields)), nil
}

// JSONFormatter formats log entries as JSON
type JSONFormatter struct{}

func (f *JSONFormatter) Format(entry *LogEntry) ([]byte, error) {
	return json.Marshal(entry)
}

// Hook interface for custom log hooks
type Hook interface {
	Fire(entry *LogEntry) error
	Levels() []LogLevel
}

// Logger represents a logger instance
type Logger struct {
	level      LogLevel
	output     io.Writer
	formatter  Formatter
	hooks      []Hook
	mu         sync.Mutex
	fields     []Field
	callerSkip int
}

// LoggerOption is a function that configures a Logger
type LoggerOption func(*Logger)

// WithLevel sets the minimum log level
func WithLevel(level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.level = level
	}
}

// WithOutput sets the output destination
func WithOutput(output io.Writer) LoggerOption {
	return func(l *Logger) {
		l.output = output
	}
}

// WithFormatter sets the log formatter
func WithFormatter(formatter Formatter) LoggerOption {
	return func(l *Logger) {
		l.formatter = formatter
	}
}

// WithHook adds a log hook
func WithHook(hook Hook) LoggerOption {
	return func(l *Logger) {
		l.hooks = append(l.hooks, hook)
	}
}

// WithFields adds default fields to the logger
func WithFields(fields ...Field) LoggerOption {
	return func(l *Logger) {
		l.fields = append(l.fields, fields...)
	}
}

// WithCallerSkip sets the number of stack frames to skip when determining the caller
func WithCallerSkip(skip int) LoggerOption {
	return func(l *Logger) {
		l.callerSkip = skip
	}
}

// NewLogger creates a new Logger instance
func NewLogger(options ...LoggerOption) *Logger {
	l := &Logger{
		level:      INFO,
		output:     os.Stdout,
		formatter:  &TextFormatter{},
		callerSkip: 2,
	}

	for _, option := range options {
		option(l)
	}

	return l
}

// log formats and writes a log message
func (l *Logger) log(level LogLevel, message string, fields ...Field) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    append(l.fields, fields...),
		Caller:    l.getCaller(),
	}

	for _, hook := range l.hooks {
		if containsLevel(hook.Levels(), level) {
			if err := hook.Fire(entry); err != nil {
				fmt.Fprintf(os.Stderr, "Error firing hook: %v\n", err)
			}
		}
	}

	formattedEntry, err := l.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting log entry: %v\n", err)
		return
	}

	_, err = l.output.Write(formattedEntry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing log entry: %v\n", err)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) getCaller() string {
	_, file, line, ok := runtime.Caller(l.callerSkip)
	if !ok {
		return "???"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func containsLevel(levels []LogLevel, level LogLevel) bool {
	for _, l := range levels {
		if l == level {
			return true
		}
	}
	return false
}

// WithField returns a new Logger with the given field added
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l.WithFields(Field{Key: key, Value: value})
}

// WithFields returns a new Logger with the given fields added
func (l *Logger) WithFields(fields ...Field) *Logger {
	newLogger := *l
	newLogger.fields = append(newLogger.fields, fields...)
	return &newLogger
}

// Trace logs a trace message
func (l *Logger) Trace(message string, fields ...Field) {
	l.log(TRACE, message, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...Field) {
	l.log(DEBUG, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...Field) {
	l.log(INFO, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...Field) {
	l.log(WARN, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...Field) {
	l.log(ERROR, message, fields...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(message string, fields ...Field) {
	l.log(FATAL, message, fields...)
}

// FileHook is a hook that writes log entries to a file
type FileHook struct {
	Writer io.Writer
}

func (h *FileHook) Fire(entry *LogEntry) error {
	formattedEntry, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = h.Writer.Write(append(formattedEntry, '\n'))
	return err
}

func (h *FileHook) Levels() []LogLevel {
	return []LogLevel{TRACE, DEBUG, INFO, WARN, ERROR, FATAL}
}

// NewFileHook creates a new FileHook with log rotation
func NewFileHook(filename string, maxSize, maxBackups, maxAge int) *FileHook {
	return &FileHook{
		Writer: &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,    // megabytes
			MaxBackups: maxBackups, // number of backups
			MaxAge:     maxAge,     // days
		},
	}
}
