package headline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takumakei/ttgo/headline"
)

func TestGet(t *testing.T) {
	tests := []struct {
		S string
		W string
	}{
		{"", ""},
		{"abc", "abc"},
		{"abc\n", "abc"},
		{"abc\ndef", "abc"},
		{"abc\ndef\n", "abc"},
		{"\nabc\ndef\n", "abc"},
		{"\n\nabc\ndef\n", "abc"},
		{"  \nabc\ndef\n", "abc"},
		{"\t\nabc\ndef\n", "abc"},
		{"  abc\ndef\n", "  abc"},
	}
	for i, e := range tests {
		got := headline.Get(e.S)
		assert.Equalf(t, e.W, got, "[%d] headline.Get(%q)", i, e.S)
	}
}
