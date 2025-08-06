package errwrap

type ErrorCode string

const (
	CodeUnknown          ErrorCode = "unknown"
	CodeInvalidArgument  ErrorCode = "invalid_argument"
	CodeNotFound         ErrorCode = "not_found"
	CodeAlreadyExists    ErrorCode = "already_exists"
	CodePermissionDenied ErrorCode = "permission_denied"
	CodeInternal         ErrorCode = "internal"
	SQLErrQuery          ErrorCode = "sql_err_query"
	SQLErrExec           ErrorCode = "sql_err_exec"
	CodeUnauthorized     ErrorCode = "unauthorized"
)
