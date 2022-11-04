package apierrors

import "errors"

var ErrLimitInvalid = errors.New("limit query parameter is no valid number")
var ErrOffsetInvalid = errors.New("offset query parameter is no valid number")
var ErrFilterInvalid = errors.New("filter query parameter is no valid number")
var ErrSortByInvalid = errors.New("sortBy query parameter is no valid number")