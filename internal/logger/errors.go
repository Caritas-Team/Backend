package logger

import (
	"fmt"
	"runtime"
)

type StackTraceError struct {
	Err        error
	Message    string
	StackTrace []string
}

func (e *StackTraceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *StackTraceError) StackTracer() []string {
	return e.StackTrace
}

func New(message string) error {
	return &StackTraceError{
		Message:    message,
		StackTrace: captureStackTrace(),
	}
}

func Errorf(format string, args ...any) error {
	return &StackTraceError{
		Message:    fmt.Sprintf(format, args...),
		StackTrace: captureStackTrace(),
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	return &StackTraceError{
		Err:        err,
		Message:    message,
		StackTrace: captureStackTrace(),
	}
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	return &StackTraceError{
		Err:        err,
		Message:    fmt.Sprintf(format, args...),
		StackTrace: captureStackTrace(),
	}
}

func captureStackTrace() []string {
	pc := make([]uintptr, 32)
	n := runtime.Callers(3, pc)
	if n == 0 {
		return []string{}
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	var stack []string
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		stackFrame := fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
		stack = append(stack, stackFrame)
	}

	return stack
}
