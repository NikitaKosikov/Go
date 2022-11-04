package service

import (
	"reflect"
	"test/internal/domain"
)

func ValidateUserField(field string) bool {
	userFields := getUserFields()
	for _, f := range userFields {
		if f == field {
			return true
		}

	}
	return false
}

func getUserFields() []string {
	var field []string
	v := reflect.ValueOf(domain.User{})
	for i := 0; i < v.Type().NumField(); i++ {
		field = append(field, v.Type().Field(i).Tag.Get("json"))
	}
	return field
}
