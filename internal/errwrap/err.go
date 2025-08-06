package errwrap

import "errors"

func WrapErrorWithReason(op string, code ErrorCode, reason string) error {
	return Wrap(op, code, errors.New(reason), map[string]any{
		"reason": reason,
	})
}

func WrapWithReason(op string, code ErrorCode, err error, reason string) error {
	return Wrap(op, code, err, map[string]any{
		"reason": reason,
	})
}

func InvalidReqError(op string, code ErrorCode, err error) error {
	return Wrap(op, code, err, map[string]any{
		"reason": "invalid_request",
	})
}
