package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/Caritas-Team/reviewer/internal/config"
)

type Logger struct {
	mu     sync.Mutex
	logger *slog.Logger
}

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	TraceIDKey   contextKey = "trace_id"
	SpanIDKey    contextKey = "span_id"
)

// Создание логгера
func NewLogger(cfg config.Config) *Logger {
	// Локальные переменные для уровня и формата логирования
	var localLevel string
	var localFormat string

	// Проверка наличия уровня и формата
	if cfg.Logging.Level != "" {
		localLevel = cfg.Logging.Level
	} else {
		slog.Warn("Отсутствует уровень логирования, используется debug")
		localLevel = "debug"
	}

	if cfg.Logging.Format != "" {
		localFormat = cfg.Logging.Format
	} else {
		slog.Warn("Отсутствует формат логирования, используется json")
		localFormat = "json"
	}

	// Проверяем корректность уровня логирования
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[localLevel] {
		slog.Warn("Некорректный уровень логирования в конфигурации, используем 'debug'.", "provided_level", localLevel)
		localLevel = "debug"
	}

	// Проверка корректности формата логирования
	validFormats := map[string]bool{"json": true, "text": true}
	if !validFormats[localFormat] {
		slog.Warn("Некорректный формат логирования в конфигурации, используем 'json'.", "provided_format", localFormat)
		localFormat = "json"
	}

	// Определяем уровень логирования
	var level slog.Level
	switch localLevel {
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	// Определяем формат логирования
	var handler slog.Handler
	if localFormat == "text" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})
	}

	return &Logger{
		logger: slog.New(handler),
	}
}

// WithFields создает новый логгер с дополнительными полями
func (l *Logger) WithFields(fields map[string]any) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	args := make([]any, 0, len(fields)*2)
	for key, value := range fields {
		args = append(args, key, value)
	}

	return &Logger{
		logger: l.logger.With(args...),
	}
}

// WithContext создает логгер с полями из контекста
func (l *Logger) WithContext(ctx context.Context) *Logger {
	fields := make(map[string]any)

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		fields["request_id"] = requestID
	}

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		fields["trace_id"] = traceID
	}

	if spanID, ok := ctx.Value(SpanIDKey).(string); ok && spanID != "" {
		fields["span_id"] = spanID
	}

	return l.WithFields(fields)
}

func (l *Logger) ErrorWithTrace(msg string, err error, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	allArgs := append([]any{"error", err}, args...)

	if stackTracer, ok := err.(interface {
		StackTrace() []string
	}); ok {
		allArgs = append(allArgs, "stack_trace", stackTracer.StackTrace())
	}

	l.logger.Error(msg, allArgs...)
}

// Методы для работы с контекстом
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, SpanIDKey, spanID)
}

// Методы логирования с мьютексом
func (l *Logger) Info(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Warn(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Debug(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Error(msg, args...)
}

// TestLog проверяет работоспособность логгера
func (l *Logger) TestLog() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Создаём временный файл для перехвата логов
	tmpfile, err := os.CreateTemp("", "log-test-*")
	if err != nil {
		return fmt.Errorf("ошибка создания временного файла: %v", err)
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name()) // Удаляем временный файл после проверки

	// Сохраняем старый вывод
	oldStdout := os.Stdout
	os.Stdout = tmpfile

	// Записываем тестовое сообщение с контекстом
	requestID := "test_request_id"
	traceID := "test_trace_id"
	ctx := context.WithValue(context.Background(), RequestIDKey, requestID)
	ctx = context.WithValue(ctx, TraceIDKey, traceID)

	msg := "This is a test log entry."
	l.WithContext(ctx).Info(msg)

	// Восстанавливаем вывод
	os.Stdout = oldStdout

	// Читаем содержимое временного файла
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return fmt.Errorf("ошибка чтения временного файла: %v", err)
	}

	// Проверяем, было ли сообщение записано
	output := string(content)
	if !strings.Contains(output, msg) ||
		!strings.Contains(output, requestID) ||
		!strings.Contains(output, traceID) {
		return fmt.Errorf("запись лока некорректна")
	}

	return nil
}
