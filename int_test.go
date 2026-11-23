// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestInt(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expectedErr string
		expected    int
	}{
		{
			name:     "no value passed",
			input:    []string{},
			expected: 0,
		},
		{
			name:     "valid value",
			input:    []string{"42"},
			expected: 42,
		},
		{
			name:     "negative value",
			input:    []string{"-7"},
			expected: -7,
		},
		{
			name:     "hex value",
			input:    []string{"0x10"},
			expected: 16,
		},
		{
			name:        "invalid value",
			input:       []string{"blabla"},
			expectedErr: `invalid argument "blabla" for "--int" flag: must be an integer`,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i int
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.IntVar(&i, "int", 0, "usage")
			err := f.Parse(repeatFlag("--int", test.input...))
			if test.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected an error; got none")
				}
				if err.Error() != test.expectedErr {
					t.Fatalf("expected error to equal %q, but was: %s", test.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error; got %q", err)
			}

			if i != test.expected {
				t.Fatalf("expected %d but got %d", test.expected, i)
			}
		})
	}
}

// TestIntOutOfRangeOn32Bit verifies that a value which does not fit in the
// platform's int size is rejected instead of being silently truncated. On
// 64-bit platforms int is 64 bits wide, so the value is in range and the test
// is skipped.
func TestIntOutOfRangeOn32Bit(t *testing.T) {
	if strconv.IntSize != 32 {
		t.Skip("int is not 32 bits on this platform")
	}

	var i int
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.IntVar(&i, "int", 0, "usage")

	// 2^31 does not fit in a 32-bit int.
	err := f.Parse([]string{"--int=2147483648"})
	if err == nil {
		t.Fatalf("expected an out-of-range error; got none (value wrapped to %d)", i)
	}
}
