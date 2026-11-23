// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestStringToString(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		flagDefault    map[string]string
		flagOpts       []zflag.Opt
		expectedErr    string
		expectedValues map[string]string
		expectedStr    string
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{},
			expectedStr:    "[]",
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: map[string]string{},
			expectedErr: `invalid argument "" for "--s2s" flag: "" must be formatted as key=value`,
		},
		{
			name:        "invalid string",
			input:       []string{"blabla"},
			flagDefault: map[string]string{},
			expectedErr: `invalid argument "blabla" for "--s2s" flag: "blabla" must be formatted as key=value`,
		},
		{
			name:           "no csv",
			input:          []string{"test=1,5"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test": "1,5"},
			expectedStr:    `[test="1,5"]`,
		},
		{
			name:           "single key value pair per arg",
			input:          []string{"test=1=1"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test": "1=1"},
			expectedStr:    `[test="1=1"]`,
		},
		{
			name:           "overrides multiple calls",
			input:          []string{"test=1", "test=5"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test": "5"},
			expectedStr:    `[test="5"]`,
		},
		{
			name:           "empty defaults",
			input:          []string{"test=1", "test2=5"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test": "1", "test2": "5"},
			expectedStr:    `[test="1" test2="5"]`,
		},
		{
			name:           "overrides default values",
			input:          []string{"test=1", "test2=5"},
			flagDefault:    map[string]string{"test2": "1", "test": "5"},
			expectedValues: map[string]string{"test": "1", "test2": "5"},
			expectedStr:    `[test="1" test2="5"]`,
		},
		{
			name:           "returns default values",
			input:          []string{},
			flagDefault:    map[string]string{"test2": "1", "test": "5"},
			expectedValues: map[string]string{"test2": "1", "test": "5"},
			expectedStr:    `[test="5" test2="1"]`,
		},
		{
			name:           "keeps whitespace",
			input:          []string{"test1=asd   ", "test2=   value", "test3=    asd   ", "test4=multi\nline\narg\npassed\nin\n"},
			flagDefault:    map[string]string{},
			expectedValues: map[string]string{"test1": "asd   ", "test2": "   value", "test3": "    asd   ", "test4": "multi\nline\narg\npassed\nin\n"},
			expectedStr:    `[test1="asd   " test2="   value" test3="    asd   " test4="multi\nline\narg\npassed\nin\n"]`,
		},
		{
			name:           "value optional",
			input:          []string{"test1"},
			flagDefault:    map[string]string{},
			flagOpts:       []zflag.Opt{zflag.OptMapValueOptional()},
			expectedValues: map[string]string{"test1": ""},
			expectedStr:    `[test1=""]`,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var s2s map[string]string
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.StringToStringVar(&s2s, "s2s", test.flagDefault, "usage", test.flagOpts...)
			err := f.Parse(repeatFlag("--s2s", test.input...))
			if test.expectedErr != "" {
				assertErr(t, err)
				assertEqualf(t, test.expectedErr, err.Error(), "expected error to equal %q, but was: %s", test.expectedErr, err)
				return
			}

			assertNoErr(t, err)

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			assertDeepEqual(t, test.expectedValues, s2s)

			int16Slice, err := f.GetStringToString("s2s")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, int16Slice)

			int16SliceGet, err := f.Get("s2s")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, int16SliceGet)

			assertEqual(t, test.expectedStr, f.Lookup("s2s").Value.String())

			defer assertNoPanic(t)()
			mustStringToString := f.MustGetStringToString("s2s")
			assertDeepEqual(t, test.expectedValues, mustStringToString)
		})
	}
}

func TestStringToStringErrors(t *testing.T) {
	t.Parallel()

	var s string
	var s2s map[string]string
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.StringToStringVar(&s2s, "s2s", map[string]string{}, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetStringToString("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetStringToString("s")
}
