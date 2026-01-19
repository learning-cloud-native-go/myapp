package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const alphaSpaceRegexString string = "^[a-zA-Z ]+$"

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

	validate.RegisterValidation("alpha_space", isAlphaSpace)

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
				msg = "This is a required field"
			case "max":
				msg = fmt.Sprintf("This must be a maximum of %s in length", err.Param())
			case "url":
				msg = "This must be a valid URL"
			case "alpha_space":
				msg = "This can only contain alphabetic and space characters"
			case "datetime":
				if err.Param() == "2006-01-02" {
					msg = "This must be a valid date"
				} else {
					msg = fmt.Sprintf("This must follow %s format", err.Param())
				}
			default:
				msg = fmt.Sprintf("something wrong on this; %s", err.Tag())
			}

			resp.Errors[err.Field()] = msg
		}

		return &resp
	}

	return nil
}

func isAlphaSpace(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(alphaSpaceRegexString)
	return reg.MatchString(fl.Field().String())
}
