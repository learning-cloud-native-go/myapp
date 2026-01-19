package validator_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"myapp/pkg/validator"
)

type testCase struct {
	name     string
	input    interface{}
	expected map[string]string
}

var tests = []*testCase{
	{
		name: `required`,
		input: struct {
			Title string `json:"title" validate:"required"`
		}{},
		expected: map[string]string{"title": "This is a required field"},
	},
	{
		name: `max`,
		input: struct {
			Course string `json:"course" validate:"max=7"`
		}{Course: "CS-0001."},
		expected: map[string]string{"course": "This must be a maximum of 7 in length"},
	},
	{
		name: `url`,
		input: struct {
			Image string `json:"image" validate:"url"`
		}{Image: "image.png"},
		expected: map[string]string{"image": "This must be a valid URL"},
	},
	{
		name: `alpha_space`,
		input: struct {
			Name string `json:"name" validate:"alpha_space"`
		}{Name: "Some Name 2"},
		expected: map[string]string{"name": "This can only contain alphabetic and space characters"},
	},
	{
		name: `date`,
		input: struct {
			Date string `json:"date" validate:"datetime=2006-01-02"`
		}{Date: "2020-02-31"},
		expected: map[string]string{"date": "This must be a valid date"},
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
			} else if !cmp.Equal(errResp.Errors, tc.expected) {
				t.Fatalf(`Expected:"%v", Got:"%v"`, tc.expected, errResp.Errors)
			}
		})
	}
}
