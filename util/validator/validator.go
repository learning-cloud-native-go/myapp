package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	alphaSpaceRegexString string = "^[a-zA-Z ]+$"
)

type ErrResponse struct {
	Errors []string `json:"errors"`
}

func New() *validator.Validate {
	validate := validator.New()
	validate.SetTagName("form")

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
			Errors: make([]string, len(fieldErrors)),
		}

		for i, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				resp.Errors[i] = fmt.Sprintf("%s is a required field", err.Field())
			case "max":
				resp.Errors[i] = fmt.Sprintf("%s must be a maximum of %s in length", err.Field(), err.Param())
			case "url":
				resp.Errors[i] = fmt.Sprintf("%s must be a valid URL", err.Field())
			case "alpha_space":
				resp.Errors[i] = fmt.Sprintf("%s can only contain alphabetic and space characters", err.Field())
			case "datetime":
				if err.Param() == "2006-01-02" {
					resp.Errors[i] = fmt.Sprintf("%s must be a valid date", err.Field())
				} else {
					resp.Errors[i] = fmt.Sprintf("%s must follow %s format", err.Field(), err.Param())
				}
			default:
				resp.Errors[i] = fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag())
			}
		}

		return &resp
	}

	return nil
}

func isAlphaSpace(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(alphaSpaceRegexString)
	return reg.MatchString(fl.Field().String())
}
