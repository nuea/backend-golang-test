package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func MapToSlice[I, O any](mapper func(I) (O, error), input []I) ([]O, error) {
	output := make([]O, 0)
	for _, i := range input {
		o, err := mapper(i)
		if err != nil {
			return nil, err
		}
		output = append(output, o)
	}
	return output, nil
}

func ValidateStruct(s interface{}) error {
	validate := validator.New(
		validator.WithRequiredStructEnabled(),
		validator.WithPrivateFieldValidation(),
	)
	if err := validate.Struct(s); err != nil {
		validationErr := err.(validator.ValidationErrors)[0]
		fieldName := validationErr.Field()
		if field, ok := reflect.TypeOf(s).Elem().FieldByName(fieldName); ok {
			fieldJSONName, ok := field.Tag.Lookup("json")
			if ok {
				fieldName = fieldJSONName
			}
		}

		return errors.New(fmt.Sprintf("%s is required.", fieldName))
	}

	return nil
}
