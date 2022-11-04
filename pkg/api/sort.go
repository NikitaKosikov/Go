package api

import (
	"fmt"
	"strings"
	apierrors "test/pkg/api/api_errors"
)

const (
	SortByParametersURL = "sortBy"
	descOrderKey           = "desc"
	ascOrderKey            = "asc"
	sortingSeparator       = ","
	fieldAndOrderSeparator = "."
	OptionsContextKey      = "sort_options"
)

type Options struct {
	Field, Order string
}

// example: sort_by=email.desc,password.asc
func ParseSort(sortBy string) ([]Options, error) {
	if sortBy == "" {
		return []Options{}, nil
	}
	options := make([]Options, 0)
	allSort := strings.Split(sortBy, sortingSeparator)
	if err := appendOptions(options, allSort); err != nil {
		return []Options{}, apierrors.ErrSortByInvalid
	}
	return options, nil
}

func appendOptions(options []Options, allSort []string) error {
	for _, s := range allSort {
		sort := strings.Split(s, fieldAndOrderSeparator)

		if len(sort) != 2 {
			return fmt.Errorf("malformed sortBy query parameter, should be field.orderdirection")
		}

		field, order := sort[0], sort[1]

		if strings.ToLower(order) != ascOrderKey && strings.ToLower(order) != descOrderKey {
			return fmt.Errorf("malformed orderdirection in sortBy query parameter, should be asc or desc")
		}

		options = append(options, Options{
			Field: field,
			Order: order,
		})
	}
	return nil
}
