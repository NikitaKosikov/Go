package api

import (
	"fmt"
	"strings"
)

const (
	FilterByParametersURL  = "filter"
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

func ParseFilters(filter string) ([]Filters, error) {
	if filter == "" {
		return []Filters{}, nil
	}
	filteres := make([]Filters, 0)
	convertedFilters := strings.Split(filter, FilterByParametersURL)
	if err := appendFilters(filteres, convertedFilters); err != nil {
		return []Filters{}, fmt.Errorf("filter query parameter is no valid number")
	}
	return filteres, nil
}

//TODO: 
func appendFilters(filters []Filters, flt []string) error {
	
	for _, f := range flt {
		fieldValue := strings.Split(f, filterFieldValueRegex)
		filter := strings.Split(f, fieldAndOrderSeparator)

		if len(filter) != 2 {
			return fmt.Errorf("malformed filter query parameter, should be field.orderdirection")
		}

		field, value := fieldValue[0], fieldValue[1]

		start := strings.Index(f, "[")
		end := strings.Index(f, "]")
		if start == -1 || end == -1 {
			return fmt.Errorf("filter parameters invalid")
		}
		operation := f[start:end]

		if err := validOperation(operation); err != nil {
			return err
		}

		filters = append(filters, Filters{
			Field:     field,
			Operation: operation,
			Value:     value,
		})
	}
	return nil
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
		return fmt.Errorf("bad operator")
	}
	return nil
}
