package validator_test

import (
	"testing"

	"myapp/util/validator"
)

type testCase struct {
	name     string
	input    interface{}
	expected string
}

var tests = []*testCase{
	{
		name: `required`,
		input: struct {
			Title string `json:"title" form:"required"`
		}{},
		expected: "title is a required field",
	},
	{
		name: `max`,
		input: struct {
			Course string `json:"course" form:"max=7"`
		}{Course: "CS-0001."},
		expected: "course must be a maximum of 7 in length",
	},
	{
		name: `url`,
		input: struct {
			Image string `json:"image" form:"url"`
		}{Image: "image.png"},
		expected: "image must be a valid URL",
	},
	{
		name: `alpha_space`,
		input: struct {
			Name string `json:"name" form:"alpha_space"`
		}{Name: "Some Name 2"},
		expected: "name can only contain alphabetic and space characters",
	},
	{
		name: `date`,
		input: struct {
			Date string `json:"date" form:"date"`
		}{Date: "2020-02-31"},
		expected: "date must be a valid date",
	},
}

func TestToErrResponse(t *testing.T) {
	vr := validator.New()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := vr.Struct(tc.input)
			if errResp := validator.ToErrResponse(err); errResp == nil || len(errResp.Errors) != 1 {
				t.Fatalf(`Expected:"{[%v]}", Got:"%v"`, tc.expected, errResp)
			} else if errResp.Errors[0] != tc.expected {
				t.Fatalf(`Expected:"%v", Got:"%v"`, tc.expected, errResp.Errors[0])
			}
		})
	}
}
