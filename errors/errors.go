package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"runtime"
)

/**
封装 github.com/pkg/errors
*/

func callers() []uintptr {
	var pcs [32]uintptr
	l := runtime.Callers(3, pcs[:])
	return pcs[:l]
}

// Error is an error with caller stack information
type Error interface {
	error
}

type item struct {
	msg   string
	stack []uintptr
}

func (i *item) Error() string {
	return i.msg
}

// Format used by go.uber.org/zap in Verbose
func (i *item) Format(s fmt.State, verb rune) {
	io.WriteString(s, i.msg)
	io.WriteString(s, "\n")

	for _, pc := range i.stack {
		fmt.Fprintf(s, "%+v\n", errors.Frame(pc))
	}
}

// New creates a new error
func New(msg string) Error {
	return &item{
		msg:   msg,
		stack: callers(),
	}
}

// Errorf creates a new error
func Errorf(format string, args ...interface{}) Error {
	return &item{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// Wrap with some extra message into err
func Wrap(err error, msg string) Error {
	if err == nil {
		return nil
	}

	e, ok := err.(*item)
	if !ok {
		return &item{
			msg:   fmt.Sprintf("%s; %s", msg, err.Error()),
			stack: callers(),
		}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// Wrapf with some extra message into err
func Wrapf(err error, format string, args ...interface{}) Error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)

	e, ok := err.(*item)
	if !ok {
		return &item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// WithStack adds caller  stack information
func WithStack(err error) Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*item); ok {
		return e
	}
	return &item{
		msg:   err.Error(),
		stack: callers(),
	}
}

func Recovery() {
	e := recover()
	if e != nil {
		s := Stack(2)
		log.Fatalf("Panic: %v\nTraceback\r:%s", e, string(s))
	}
}

func RecoverStackWithoutLF() {
	e := recover()
	if e != nil {
		s := StackWithoutLF(3)
		log.Fatalf("Panic: %v Traceback:%s", e, s)
	}
}
