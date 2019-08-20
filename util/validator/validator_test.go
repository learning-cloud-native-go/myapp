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
