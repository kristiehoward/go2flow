package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTagInfo(t *testing.T) {
	type testCase struct {
		Tag         string
		Name        string
		IsOptional  bool
		Description string
	}

	testCases := []testCase{
		{
			`json:"screenshots"`,
			"screenshots",
			false,
			"Standard json tag",
		},
		{
			`json:"data_2_take,omitempty"`,
			"data_2_take",
			true,
			"Standard optional json tag",
		},
		{
			`json:",omitempty"`,
			"",
			false,
			"Missing optional json tag",
		},
		{
			`json:""`,
			"",
			false,
			"Missing json tag",
		},
		{
			`kub:"otherthings"`,
			"",
			false,
			"Non-json tag",
		},
		{
			`kub:"otherthings,omitempty"`,
			"",
			false,
			"Non-json tag with optional",
		},
	}

	var name string
	var isOptional bool

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			name, isOptional = GetTagInfo(tc.Tag)
			assert.Equal(t, tc.Name, name)
			if tc.IsOptional {
				assert.True(t, isOptional)
			} else {
				assert.False(t, isOptional)
			}
		})
	}

}
