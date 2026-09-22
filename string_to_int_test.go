// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    map[string]int
		input          []string
		expectedErr    string
		expectedValues map[string]int
		expectedStr    string
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    map[string]int{},
			expectedErr:    "",
			expectedValues: map[string]int{},
			expectedStr:    "[]",
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: map[string]int{},
			expectedErr: `invalid argument "" for "--s2i" flag:  must be formatted as key=value`,
		},
		{
			name:        "invalid int",
			input:       []string{"blabla"},
			flagDefault: map[string]int{},
			expectedErr: `invalid argument "blabla" for "--s2i" flag: blabla must be formatted as key=value`,
		},
		{
			name:        "no csv",
			input:       []string{"test=1,5"},
			flagDefault: map[string]int{},
			expectedErr: `invalid argument "test=1,5" for "--s2i" flag: strconv.Atoi: parsing "1,5": invalid syntax`,
		},
		{
			name:        "single key value pair per arg",
			input:       []string{"test=1=1"},
			flagDefault: map[string]int{},
			expectedErr: `invalid argument "test=1=1" for "--s2i" flag: strconv.Atoi: parsing "1=1": invalid syntax`,
		},
		{
			name:           "overrides multiple calls",
			input:          []string{"test=1", "test=5"},
			flagDefault:    map[string]int{},
			expectedValues: map[string]int{"test": 5},
			expectedStr:    "[test=5]",
		},
		{
			name:           "empty defaults",
			input:          []string{"test=1", "test2=5"},
			flagDefault:    map[string]int{},
			expectedValues: map[string]int{"test": 1, "test2": 5},
			expectedStr:    "[test=1 test2=5]",
		},
		{
			name:           "overrides default values",
			input:          []string{"test=1", "test2=5"},
			flagDefault:    map[string]int{"test2": 1, "test": 5},
			expectedValues: map[string]int{"test": 1, "test2": 5},
			expectedStr:    "[test=1 test2=5]",
		},
		{
			name:           "returns default values",
			input:          []string{},
			flagDefault:    map[string]int{"test2": 1, "test": 5},
			expectedValues: map[string]int{"test2": 1, "test": 5},
			expectedStr:    "[test=5 test2=1]",
		},
		{
			name:           "trims input",
			input:          []string{"test=    1", "test2=5     ", "test3=     9     "},
			flagDefault:    map[string]int{},
			expectedValues: map[string]int{"test": 1, "test2": 5, "test3": 9},
			expectedStr:    "[test=1 test2=5 test3=9]",
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var s2i map[string]int
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.StringToIntVar(&s2i, "s2i", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--s2i", test.input...))
			if test.expectedErr != "" {
				assertErr(t, err)
				assertEqualf(t, test.expectedErr, err.Error(), "expected error to equal %q, but was: %s", test.expectedErr, err)
				return
			}

			assertNoErr(t, err)

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			assertDeepEqual(t, test.expectedValues, s2i)

			s2iGetS2I, err := f.GetStringToInt("s2i")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, s2iGetS2I)

			s2iGet, err := f.Get("s2i")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, s2iGet)

			assertEqual(t, test.expectedStr, f.Lookup("s2i").Value.String())

			defer assertNoPanic(t)()
			mustStringToInt := f.MustGetStringToInt("s2i")
			assertDeepEqual(t, test.expectedValues, mustStringToInt)
		})
	}
}

func TestStringToIntErrors(t *testing.T) {
	t.Parallel()

	var s string
	var s2i map[string]int
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.StringToIntVar(&s2i, "s2i", map[string]int{}, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetStringToInt("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetStringToInt("s")
}
