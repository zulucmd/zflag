// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"errors"
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
