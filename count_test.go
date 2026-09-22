// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		expectedErr   string
		expectedValue int
	}{
		{
			name:          "no flags",
			input:         []string{},
			expectedValue: 0,
		},
		{
			name:          "single count",
			input:         []string{"-v"},
			expectedValue: 1,
		},
		{
			name:          "multiple times",
			input:         []string{"-vvv"},
			expectedValue: 3,
		},
		{
			name:          "multiple times separated",
			input:         []string{"-v", "-v", "-v"},
			expectedValue: 3,
		},
		{
			name:          "multiple times interchanged and separated",
			input:         []string{"-v", "--verbose", "-v"},
			expectedValue: 3,
		},
		{
			name:          "multiple times with value",
			input:         []string{"-v=3", "-v"},
			expectedValue: 4,
		},
		{
			name:          "long opt with value",
			input:         []string{"--verbose=0"},
			expectedValue: 0,
		},
		{
			name:          "single with value",
			input:         []string{"-v=0"},
			expectedValue: 0,
		},
		{
			name:          "",
			input:         []string{"-v=a"},
			expectedErr:   `invalid argument "a" for "-v, --verbose" flag: must be an integer`,
			expectedValue: 0,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var verbose int
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.CountVar(&verbose, "verbose", "usage", zflag.OptShorthand('v'))
			err := f.Parse(test.input)
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, verbose)

			getVerbose, err := f.GetCount("verbose")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, getVerbose)

			getVerboseGet, err := f.Get("verbose")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, getVerboseGet)

			defer assertNoPanic(t)()
			mustBool := f.MustGetCount("verbose")
			assertEqual(t, test.expectedValue, mustBool)
		})
	}
}

func TestCountErrors(t *testing.T) {
	t.Parallel()

	var s string
	var count int
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.CountVar(&count, "count", "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetBool("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetCount("s")
}
