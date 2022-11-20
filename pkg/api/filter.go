package api

import (
	"fmt"
	"regexp"
	"strings"
	apierrors "test/pkg/api/api_errors"
)

const (
	FilterByParametersURL  = "filter"
	FiltersSeparator       = ","
	filterFieldValueRegex  = "\\[.+\\]="
	opearatorEqual         = "eq"
	opearatorNotEqual      = "ne"
	opearatorGreaterThan   = "gt"
	opearatorGreaterThanEq = "gte"
	opearatorLowerThan     = "lt"
	opearatorLowerThanEq   = "lte"
	opearatorAnd           = "and"
	opearatorOr            = "or"
)

type Filters struct {
	Field, Operation, Value string
}

// example: filter=email[eq]=email,password[eq]=password
func ParseFilters(filter string) ([]Filters, error) {
	if filter == "" {
		return []Filters{}, nil
	}

	convertedFilters := strings.Split(filter, FiltersSeparator)
	filteres, err := appendFilters(convertedFilters)
	if err != nil {
		return []Filters{}, err
	}
	return filteres, nil
}

// TODO:
func appendFilters(flt []string) ([]Filters, error) {
	re := regexp.MustCompile("\\[.+\\]=")
	filters:=make([]Filters, 0)
	for _, f := range flt {

		fieldValue := re.Split(f, -1)

		if len(fieldValue) != 2 {
			return []Filters{}, apierrors.ErrFilterInvalid
		}

		field, value := fieldValue[0], fieldValue[1]

		start := strings.Index(f, "[")
		end := strings.Index(f, "]")
		if start == -1 || end == -1 {
			return []Filters{}, apierrors.ErrFilterInvalid
		}
		operation := f[start+1 : end]

		if err := validOperation(operation); err != nil {
			return []Filters{}, err
		}

		filters = append(filters, Filters{
			Field:     field,
			Operation: operation,
			Value:     value,
		})
	}
	return filters, nil
}

func validOperation(operation string) error {
	switch operation {
	case opearatorEqual:
	case opearatorNotEqual:
	case opearatorGreaterThan:
	case opearatorGreaterThanEq:
	case opearatorLowerThan:
	case opearatorLowerThanEq:
	case opearatorAnd:
	case opearatorOr:
	default:
		return apierrors.ErrFilterOperatorInvalid
	}
	return nil
}
