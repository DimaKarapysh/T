package delivery

type UserError struct {
	error
	code string
}

func NewUserError(code string, err error) UserError {
	return UserError{
		error: err,
		code:  code,
	}
}
