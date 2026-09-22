// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestUint(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expectedErr string
		expected    uint
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
			name:     "hex value",
			input:    []string{"0x10"},
			expected: 16,
		},
		{
			name:        "negative value",
			input:       []string{"-7"},
			expectedErr: `invalid argument "-7" for "--uint" flag: must be a non-negative integer`,
		},
		{
			name:        "invalid value",
			input:       []string{"blabla"},
			expectedErr: `invalid argument "blabla" for "--uint" flag: must be a non-negative integer`,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var u uint
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.UintVar(&u, "uint", 0, "usage")
			err := f.Parse(repeatFlag("--uint", test.input...))
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

			if u != test.expected {
				t.Fatalf("expected %d but got %d", test.expected, u)
			}
		})
	}
}

// TestUintOutOfRangeOn32Bit verifies that a value which does not fit in the
// platform's uint size is rejected instead of being silently truncated. On
// 64-bit platforms uint is 64 bits wide, so the value is in range and the test
// is skipped.
func TestUintOutOfRangeOn32Bit(t *testing.T) {
	if strconv.IntSize != 32 {
		t.Skip("uint is not 32 bits on this platform")
	}

	var u uint
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.UintVar(&u, "uint", 0, "usage")

	// 2^32 does not fit in a 32-bit uint.
	err := f.Parse([]string{"--uint=4294967296"})
	if err == nil {
		t.Fatalf("expected an out-of-range error; got none (value wrapped to %d)", u)
	}
}
