package gvalidator

import (
	"fmt"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var instance = validator.New(validator.WithRequiredStructEnabled())
var customTagsRegistered = false

func Instance() *validator.Validate {
	err := RegisterCustomTags()
	if err != nil {
		panic(fmt.Errorf("failed to register custom tags: %w", err))
	}

	return instance
}

func RegisterCustomTags() error {
	if customTagsRegistered {
		return nil
	}

	if err := registerPortTag(); err != nil {
		return fmt.Errorf("failed to register port tag: %w", err)
	}

	if err := instance.RegisterValidation("notblank", validators.NotBlank); err != nil {
		return fmt.Errorf("failed to register notblank tag %w", err)
	}

	customTagsRegistered = true
	return nil
}

func registerPortTag() error {
	return instance.RegisterValidation("port", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		var v uint64
		switch field.Kind() {
		case reflect.String:
			parsedV, err := strconv.ParseUint(field.String(), 10, 64)
			if err != nil {
				panic(err)
			}

			v = parsedV
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v = uint64(field.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v = field.Uint()
		default:
			panic(fmt.Sprintf("Bad field type %T", field.Interface()))
		}

		return v >= 1 && v <= 65535
	})
}
