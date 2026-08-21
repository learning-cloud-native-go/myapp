package validator_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"myapp/pkg/validator"
)

type testCase struct {
	name     string
	input    any
	expected map[string]string
}

var tests = []*testCase{
	{
		name: `required`,
		input: struct {
			Title string `json:"title" validate:"required"`
		}{},
		expected: map[string]string{"title": "This field is required"},
	},
	{
		name: `max`,
		input: struct {
			Course string `json:"course" validate:"max=7"`
		}{Course: "CS-0001."},
		expected: map[string]string{"course": "Must be no more than 7 characters"},
	},
	{
		name: `url`,
		input: struct {
			Image string `json:"image" validate:"url"`
		}{Image: "image.png"},
		expected: map[string]string{"image": "Must be a valid URL"},
	},
	{
		name: `oneof`,
		input: struct {
			Status string `json:"status" validate:"required,oneof=pending verified"`
		}{
			Status: "invalid",
		},
		expected: map[string]string{"status": "Must be one of pending, verified"},
	},
	{
		name: `date`,
		input: struct {
			Date string `json:"date" validate:"datetime=2006-01-02"`
		}{Date: "2020-02-31"},
		expected: map[string]string{"date": "Must be a valid date"},
	},
}

func TestToErrResponse(t *testing.T) {
	vr := validator.New()

	for _, tc := range tests {
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
