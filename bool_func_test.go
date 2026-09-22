// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestBoolFunc(t *testing.T) {
	var count int
	fn := func(_ string) error {
		count++
		return nil
	}

	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.BoolFunc("func", "Callback function", fn)

	assertNoErr(t, f.Parse([]string{"--func", "--func=1", "--func=false"}))
	assertEqual(t, 3, count)
}

func TestBoolFuncShorthand(t *testing.T) {
	var count int
	fn := func(_ string) error {
		count++
		return nil
	}

	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.BoolFunc("bfunc", "Callback function", fn, zflag.OptShorthand('b'))

	assertNoErr(t, f.Parse([]string{"--bfunc", "--bfunc=0", "--bfunc=false", "-b", "-b=0"}))
	assertEqual(t, 5, count)
}

func TestBoolFuncError(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.BoolFunc("func", "usage", func(string) error { return errors.New("callback error") })

	assertErrMsg(t, `invalid argument "true" for "--func" flag: callback error`, f.Parse([]string{"--func"}))
}

func TestBoolFuncUsage(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		usage    string
		expected string
	}{
		{
			name:     "regular boolfunc flag",
			flagName: "flag1",
			usage:    "usage message",
			expected: "--flag1   usage message",
		},
		{
			name:     "boolfunc flag with placeholder name",
			flagName: "flag2",
			usage:    "usage message with `name` placeholder",
			expected: "--flag2 name   usage message with name placeholder",
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			f := zflag.NewFlagSet("unittest", zflag.ContinueOnError)
			f.BoolFunc(test.flagName, test.usage, func(string) error { return nil })

			assertEqual(t, test.expected, strings.TrimSpace(f.FlagUsagesWrapped(80)))
		})
	}
}
