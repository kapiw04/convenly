package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
)

const logFilename = "convenly"

type CustomHandler struct {
	writer io.Writer
	isFile bool
	level  slog.Level
	attrs  []slog.Attr
	group  string
}

func NewCustomHandler(writer io.Writer, isFile bool, level slog.Level) *CustomHandler {
	return &CustomHandler{
		writer: writer,
		isFile: isFile,
		level:  level,
	}
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	timestamp := record.Time.Format("15:04:05.000")

	levelStr := record.Level.String()
	var coloredLevel string
	switch record.Level {
	case slog.LevelDebug:
		coloredLevel = color.MagentaString(levelStr)
	case slog.LevelInfo:
		coloredLevel = color.BlueString(levelStr)
	case slog.LevelWarn:
		coloredLevel = color.YellowString(levelStr)
	case slog.LevelError:
		coloredLevel = color.RedString(levelStr)
	default:
		coloredLevel = levelStr
	}

	message := record.Message
	if len(h.attrs) > 0 {
		var attrs []string
		for _, attr := range h.attrs {
			attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
		}
		message += " " + strings.Join(attrs, " ")
	}

	record.Attrs(func(attr slog.Attr) bool {
		message += fmt.Sprintf(" %s=%v", attr.Key, attr.Value)
		return true
	})

	if h.isFile {
		_, err := fmt.Fprintf(h.writer, "[%s] %s %s\n", levelStr, timestamp, message)
		return err
	} else {
		_, err := fmt.Fprintf(h.writer, "%-8s %-12s %s\n", "["+coloredLevel+"]", color.WhiteString(timestamp), color.CyanString(message))
		return err
	}
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs))
	copy(newAttrs, h.attrs)
	newAttrs = append(newAttrs, attrs...)
	return &CustomHandler{
		writer: h.writer,
		isFile: h.isFile,
		level:  h.level,
		attrs:  newAttrs,
		group:  h.group,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		writer: h.writer,
		isFile: h.isFile,
		level:  h.level,
		attrs:  h.attrs,
		group:  name,
	}
}

func InitializeLogger(logDir string) {
	loggerFile := getLoggerFile(logDir)
	debug.SetCrashOutput(loggerFile, debug.CrashOptions{})

	fileHandler := NewCustomHandler(loggerFile, true, slog.LevelDebug)
	consoleHandler := NewCustomHandler(os.Stdout, false, slog.LevelDebug)

	handler := &MultiHandler{handlers: []slog.Handler{fileHandler, consoleHandler}}

	log := slog.New(handler)
	slog.SetDefault(log)

	slog.Info("Logger successfully initialized")
}

type MultiHandler struct {
	handlers []slog.Handler
}

func (mh *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range mh.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (mh *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, h := range mh.handlers {
		if h.Enabled(ctx, record.Level) {
			if err := h.Handle(ctx, record); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mh *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (mh *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

func getLoggerFile(logDir string) *os.File {
	currentFilePath := filepath.Join(logDir, fmt.Sprintf("%s_current.log", logFilename))
	previousFilePath := filepath.Join(logDir, fmt.Sprintf("%s_previous.log", logFilename))

	if _, err := os.Stat(currentFilePath); err == nil {
		os.Remove(previousFilePath)
		err := os.Rename(currentFilePath, previousFilePath)
		if err != nil {
			fmt.Println("Couldn't rename current log file to previous")
		}
	}

	f, err := os.OpenFile(currentFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Couldn't open file for logging")
	}
	return f
}
