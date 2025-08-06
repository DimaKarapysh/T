package errwrap

import (
	"errors"
	"fmt"
)

// ErrorInfo is a structured errwrap
type ErrorInfo struct {
	Op   string
	Code ErrorCode
	Err  error
	Meta map[string]any
}

// Error returns errwrap message
func (e *ErrorInfo) Error() string {
	if e.Meta != nil {
		return fmt.Sprintf("%s: [%s] %v | meta: %v", e.Op, e.Code, e.Err, e.Meta)
	}
	return fmt.Sprintf("%s: [%s] %v", e.Op, e.Code, e.Err)
}

func (e *ErrorInfo) Unwrap() error {
	return e.Err
}

func CodeOf(err error) ErrorCode {
	var ei *ErrorInfo
	if errors.As(err, &ei) {
		return ei.Code
	}
	return CodeUnknown
}

// --- Helpers ---

func Wrap(op string, code ErrorCode, err error, meta map[string]any) error {
	if err == nil {
		err = errors.New("nil errwrap")
	}
	return &ErrorInfo{
		Op:   op,
		Code: code,
		Err:  err,
		Meta: meta,
	}
}

// Common short-hand functions

func InvalidArg(op string, err error) error {
	return Wrap(op, CodeInvalidArgument, err, nil)
}

func NotFound(op string, err error) error {
	return Wrap(op, CodeNotFound, err, nil)
}

func AlreadyExists(op string, err error) error {
	return Wrap(op, CodeAlreadyExists, err, nil)
}

func PermissionDenied(op string, err error) error {
	return Wrap(op, CodePermissionDenied, err, nil)
}

func Internal(op string, err error) error {
	return Wrap(op, CodeInternal, err, nil)
}

func ErrUserIDEmpty(op string) error {
	return Wrap(op, CodeInvalidArgument, errors.New("user_id_empty"), nil)
}

func ErrSessionIDEmpty(op string) error {
	return Wrap(op, CodeInvalidArgument, errors.New("session_id_empty"), nil)
}

func SQLQueryError(op string, err error) error {
	return Wrap(op, SQLErrQuery, err, nil)
}

func SQLExecError(op string, err error) error {
	return Wrap(op, SQLErrExec, err, nil)
}
