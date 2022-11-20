package apierrors

import "errors"

var ErrLimitInvalid = NewApiErr("limit query parameter is no valid number")
var ErrOffsetInvalid = NewApiErr("offset query parameter is no valid number")
var ErrFilterInvalid = NewApiErr("malformed filter query parameter, should be field[operator]=value")
var ErrFilterOperatorInvalid = NewApiErr("invalid filter operator")
var ErrSortByInvalid = NewApiErr("sortBy query parameter is no valid number")

type ApiError struct {
	Err error
}


func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) Unwrap() error {
	return e.Err
}


func NewApiErr(message string) *ApiError {
	return &ApiError{
		Err: errors.New(message),
	}
}
