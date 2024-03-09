package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type EchoValidator struct {
	Validator *validator.Validate
}

func (v *EchoValidator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		var errMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			jsonTag := GetTagName(i, err.Field())
			if err.Tag() == "required" {
				errMessages = append(errMessages, fmt.Sprintf("Field %s is %s", jsonTag, err.Tag()))
				continue
			}
			errMessages = append(errMessages, fmt.Sprintf("Field '%s' is invalid", jsonTag))
		}
		return errors.New(strings.Join(errMessages, "; "))
	}

	return nil
}

func GetTagName(i interface{}, field string) string {
	f, ok := reflect.TypeOf(i).Elem().FieldByName(field)
	if ok {
		json, _ := f.Tag.Lookup("json")
		return json
	}
	return ""
}
