package handler_test

import (
	"testing"

	"github.com/crowmw/risiti/handler"
	"github.com/stretchr/testify/assert"
)

var testcases = []struct {
	name     string
	expected string
	input    string
}{
	{"replace special charactes with -", "file-name-2", "file!@#$name@#$%2"},
	{"trim spaces", "filename2", " filename2 "},
	{"replace spaces with -", "filename-2-1", "filename 2 1"},
	{"lowercase input", "filename", "FiLeNaMe"},
	{"empty stirng", "", ""},
	{"only special char", "", "$%^"},
}

func TestCreateSlug(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			got := handler.CreateSlug(tc.input)

			assert.Equal(got, tc.expected)
		})
	}
}
