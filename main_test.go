package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO Kristie 10/24/17 - Update to include the edge cases in
// https://golang.org/pkg/encoding/json/#Marshal
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
			`protobuf:"otherthings"`,
			"",
			false,
			"Non-json tag",
		},
		{
			`protobuf:"otherthings,omitempty"`,
			"",
			false,
			"Non-json tag with optional",
		},
		{
			`json:"date,omitempty" protobuf:"bytes,1,opt,name=name"`,
			"date",
			true,
			"json tag first with additional defns",
		},
		{
			`protobuf:"bytes,1,opt,name=name" json:"date,omitempty"`,
			"date",
			true,
			"json tag second with additional defns",
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
