package camino

import (
	"strconv"
	"testing"

	"github.com/boggydigital/testo"
)

type sanInOut struct {
	in     string
	expOut string
}

func TestSanitize(t *testing.T) {

	tests := []sanInOut{
		{"", ""},
		{"abc", "abc"},
		{"\n", safeStr},
		{"abc\n", "abc" + safeStr},
		{"\nabc", safeStr + "abc"},
	}

	for _, ch := range problematicChars {
		tests = append(tests, sanInOut{string(ch), safeStr})
	}
	for _, ch := range asciiControlChars {
		tests = append(tests, sanInOut{string(ch), safeStr})
	}
	for _, name := range ntfsIntFilenames {
		tests = append(tests, sanInOut{name, safeStr + name})
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii+1), func(t *testing.T) {
			testo.EqualValues(t, Sanitize(tt.in), tt.expOut)
		})
	}
}
