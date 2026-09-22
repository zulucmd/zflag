// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestC128Slice(t *testing.T) {
	tests := []struct {
		name              string
		flagDefault       []complex128
		input             []string
		expectedErr       string
		expectedValues    []complex128
		visitor           func(f *zflag.Flag)
		expectedStrValues string
		expectedGetSlice  []string
	}{
		{
			name:              "no value passed",
			input:             []string{},
			flagDefault:       []complex128{},
			expectedErr:       "",
			expectedValues:    []complex128{},
			expectedStrValues: "[]",
			expectedGetSlice:  []string{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []complex128{},
			expectedErr: `invalid argument "" for "--c128s" flag: must be a complex number`,
		},
		{
			name:        "invalid c128s",
			input:       []string{"blabla"},
			flagDefault: []complex128{},
			expectedErr: `invalid argument "blabla" for "--c128s" flag: must be a complex number`,
		},
		{
			name:        "no csv",
			input:       []string{"1.0,2.0"},
			flagDefault: []complex128{},
			expectedErr: `invalid argument "1.0,2.0" for "--c128s" flag: must be a complex number`,
		},
		{
			name:              "multiple values passed",
			input:             []string{"1.0", "2.0"},
			flagDefault:       []complex128{},
			expectedValues:    []complex128{1.0, 2.0},
			expectedStrValues: "[(1.000000+0.000000i) (2.000000+0.000000i)]",
			expectedGetSlice:  []string{"(1.000000+0.000000i)", "(2.000000+0.000000i)"},
		},
		{
			name:              "with default values",
			input:             []string{"1.0", "2.0"},
			flagDefault:       []complex128{2.0, 1.0},
			expectedValues:    []complex128{1.0, 2.0},
			expectedStrValues: "[(1.000000+0.000000i) (2.000000+0.000000i)]",
			expectedGetSlice:  []string{"(1.000000+0.000000i)", "(2.000000+0.000000i)"},
		},
		{
			name:  "replace values",
			input: []string{"1.0"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"0+2i"})
				}
			},
			expectedValues:    []complex128{complex(0, 2)},
			expectedStrValues: "[(0.000000+2.000000i)]",
			expectedGetSlice:  []string{"(0.000000+2.000000i)"},
		},
		{
			name:  "replace values error",
			input: []string{"1.0"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					err := val.Replace([]string{"notc128"})
					assertErr(t, err)
				}
			},
			expectedValues:    []complex128{complex(1, 0)},
			expectedStrValues: "[(1.000000+0.000000i)]",
			expectedGetSlice:  []string{"(1.000000+0.000000i)"},
		},
		{
			name:  "add values",
			input: []string{"1.0"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Append("2.0")
				}
			},
			expectedValues:    []complex128{complex(1, 0), complex(2, 0)},
			expectedStrValues: "[(1.000000+0.000000i) (2.000000+0.000000i)]",
			expectedGetSlice:  []string{"(1.000000+0.000000i)", "(2.000000+0.000000i)"},
		},
		{
			name:  "add values error",
			input: []string{"1.0"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					err := val.Append("asd")
					if err == nil {
						t.Errorf("Expected an error when appending, got %s", err)
					}
				}
			},
			flagDefault:       nil,
			expectedValues:    []complex128{complex(1, 0)},
			expectedStrValues: "[(1.000000+0.000000i)]",
			expectedGetSlice:  []string{"(1.000000+0.000000i)"},
		},
		{
			name:              "nil default",
			input:             []string{},
			flagDefault:       nil,
			expectedValues:    nil,
			expectedStrValues: "[]",
			expectedGetSlice:  []string{},
		},
		{
			name:              "valid c128s",
			input:             []string{"1.0", "2.0", "3.0", "0+2i", "1", "2i", "2.5+3.1i"},
			expectedValues:    []complex128{1.0, 2.0, 3.0, complex(0, 2), complex(1, 0), complex(0, 2), complex(2.5, 3.1)},
			expectedStrValues: "[(1.000000+0.000000i) (2.000000+0.000000i) (3.000000+0.000000i) (0.000000+2.000000i) (1.000000+0.000000i) (0.000000+2.000000i) (2.500000+3.100000i)]",
			expectedGetSlice:  []string{"(1.000000+0.000000i)", "(2.000000+0.000000i)", "(3.000000+0.000000i)", "(0.000000+2.000000i)", "(1.000000+0.000000i)", "(0.000000+2.000000i)", "(2.500000+3.100000i)"},
		},
		{
			name:              "trims input",
			input:             []string{" 1.0 ", "   2.0", "3.0   ", "  0+2i", "1"},
			expectedValues:    []complex128{1.0, 2.0, 3.0, complex(0, 2), complex(1, 0)},
			expectedStrValues: "[(1.000000+0.000000i) (2.000000+0.000000i) (3.000000+0.000000i) (0.000000+2.000000i) (1.000000+0.000000i)]",
			expectedGetSlice:  []string{"(1.000000+0.000000i)", "(2.000000+0.000000i)", "(3.000000+0.000000i)", "(0.000000+2.000000i)", "(1.000000+0.000000i)"},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var c128s []complex128
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.Complex128SliceVar(&c128s, "c128s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--c128s", test.input...))
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			assertDeepEqual(t, test.expectedValues, c128s)

			getC128s, err := f.GetComplex128Slice("c128s")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, getC128s)

			getC128sGet, err := f.Get("c128s")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, getC128sGet)

			flag := f.Lookup("c128s")
			assertEqual(t, test.expectedStrValues, flag.Value.String())

			sliced := flag.Value.(zflag.SliceValue)
			assertDeepEqual(t, test.expectedGetSlice, sliced.GetSlice())

			defer assertNoPanic(t)()
			mustComplex128Slice := f.MustGetComplex128Slice("c128s")
			assertDeepEqual(t, test.expectedValues, mustComplex128Slice)
		})
	}
}

func TestComplex128SliceErrors(t *testing.T) {
	var s string
	var c128s []complex128
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.Complex128SliceVar(&c128s, "c128s", []complex128{}, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetComplex128Slice("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetComplex128Slice("s")
}
