// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestComplex128(t *testing.T) {
	tests := []struct {
		name          string
		flagDefault   complex128
		input         []string
		expectedErr   string
		expectedValue complex128
		extraOpts     []zflag.Opt
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   complex(1, 0),
			expectedErr:   "",
			expectedValue: complex(1, 0),
		},
		{
			name:        "empty value passed",
			input:       repeatFlag("--c128", ""),
			flagDefault: complex(1, 0),
			expectedErr: `invalid argument "" for "--c128" flag: must be a complex number`,
		},
		{
			name:        "invalid complex128",
			input:       repeatFlag("--c128", "blabla"),
			flagDefault: complex(1, 0),
			expectedErr: `invalid argument "blabla" for "--c128" flag: must be a complex number`,
		},
		{
			name:        "no csv",
			input:       repeatFlag("--c128", "1.0,1.0"),
			flagDefault: complex(0, 0),
			expectedErr: `invalid argument "1.0,1.0" for "--c128" flag: must be a complex number`,
		},
		{
			name:          "accepts separate value without no-",
			input:         []string{"--c128", "1.0"},
			flagDefault:   complex(0, 0),
			expectedValue: complex(1, 0),
		},
		{
			name:          "repeated value",
			input:         repeatFlag("--c128", "1.0", "3.0"),
			flagDefault:   complex(0, 0),
			expectedValue: complex(3, 0),
		},
		{
			name:          "with default values",
			input:         []string{},
			flagDefault:   complex(4, 0),
			expectedValue: complex(4, 0),
		},
		{
			name:          "trims input",
			input:         repeatFlag("--c128", " 1.0 "),
			expectedValue: complex(1, 0),
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var c128 complex128
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Complex128Var(&c128, "c128", test.flagDefault, "usage", test.extraOpts...)
			err := f.Parse(test.input)
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, c128)

			getBS, err := f.GetComplex128("c128")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, getBS)

			getBSGet, err := f.Get("c128")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, getBSGet)

			defer assertNoPanic(t)()
			mustComplex128 := f.MustGetComplex128("c128")
			assertEqual(t, test.expectedValue, mustComplex128)
		})
	}
}

func TestComplex128Errors(t *testing.T) {
	var s string
	var c128 complex128
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.Complex128Var(&c128, "c128", complex(1, 0), "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetComplex128("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetComplex128("s")
}
