package api

import (
	"regexp"
	"strings"
	apierrors "test/pkg/api/api_errors"
)

const (
	FilterByParametersURL  = "filter"
	FilterPattern          = "\\[.+\\]="
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

func appendFilters(flt []string) ([]Filters, error) {
	filters := make([]Filters, 0)
	for _, f := range flt {
		field, value, err := extractFieldAndValue(f)
		if err != nil {
			return []Filters{}, err
		}

		operation, err := extractOperation(f)
		if err != nil {
			return []Filters{}, err
		}

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

func extractFieldAndValue(filter string) (string, string, error) {
	filterRegex := regexp.MustCompile(FilterPattern)
	fieldValue := filterRegex.Split(filter, -1)

	if len(fieldValue) != 2 {
		return "", "", apierrors.ErrFilterInvalid
	}

	field, value := fieldValue[0], fieldValue[1]
	return field, value, nil
}

func extractOperation(filter string) (string, error) {
	start := strings.Index(filter, "[")
	end := strings.Index(filter, "]")
	if start == -1 || end == -1 {
		return "", apierrors.ErrFilterInvalid
	}
	return filter[start+1 : end], nil
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
