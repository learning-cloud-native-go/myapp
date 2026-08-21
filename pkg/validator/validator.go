package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Errors map[string]string `json:"errors"`
}

func New() *validator.Validate {
	validate := validator.New()

	// Using the names which have been specified for JSON representations of structs, rather than normal Go field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// custom validators

	return validate
}

func ToErrResponse(err error) *ErrResponse {
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		resp := ErrResponse{
			Errors: make(map[string]string, len(fieldErrors)),
		}

		for _, err := range fieldErrors {
			var msg string

			switch err.Tag() {
			case "required":
				msg = "This field is required"
			case "max":
				msg = fmt.Sprintf("Must be no more than %s characters", err.Param())
			case "url":
				msg = "Must be a valid URL"
			case "oneof":
				msg = fmt.Sprintf("Must be one of %s", strings.ReplaceAll(err.Param(), " ", ", "))
			case "datetime":
				if err.Param() == "2006-01-02" {
					msg = "Must be a valid date"
				} else {
					msg = fmt.Sprintf("Must follow %s format", err.Param())
				}
			default:
				msg = "This value is invalid"
			}

			resp.Errors[err.Field()] = msg
		}

		return &resp
	}

	return nil
}
