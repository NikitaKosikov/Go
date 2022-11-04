package api

import (
	"strconv"
	apierrors "test/pkg/api/api_errors"
)

const (
	LimitByParametersURL  = "limit"
	OffsetByParametersURL = "offset"
)

type Pagination struct {
	Limit, Offset int64
}

func NewPagination(limit, offset string) (Pagination, error) {

	convertedLimit, err := convertLimit(limit)
	if err != nil {
		return Pagination{}, err
	}
	convertedOffset, err := convertOffset(offset)
	if err != nil {
		return Pagination{}, err
	}
	return Pagination{
		Limit:  convertedLimit,
		Offset: convertedOffset,
	}, nil
}

func convertLimit(limit string) (int64, error) {

	if limit != "" {
		limit, err := strconv.Atoi(limit)
		if err != nil || limit < 0 {
			return 0, apierrors.ErrLimitInvalid
		}
		return int64(limit), nil
	}
	return 0, nil
}

func convertOffset(offset string) (int64, error) {

	if offset != "" {
		offset, err := strconv.Atoi(offset)
		if err != nil || offset < 0 {
			return 0, apierrors.ErrOffsetInvalid
		}
		return int64(offset), nil
	}
	return 0, nil
}
