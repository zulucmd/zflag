// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestInt8Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []int8
		input          []string
		expectedErr    string
		expectedValues []int8
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []int8{},
			expectedErr:    "",
			expectedValues: []int8{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []int8{},
			expectedErr: `invalid argument "" for "--i8s" flag: must be an integer`,
		},
		{
			name:        "invalid int8",
			input:       []string{"blabla"},
			flagDefault: []int8{},
			expectedErr: `invalid argument "blabla" for "--i8s" flag: must be an integer`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []int8{},
			expectedErr: `invalid argument "1,5" for "--i8s" flag: must be an integer`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []int8{},
			expectedValues: []int8{1, 5},
		},
		{
			name:           "with default values",
			input:          []string{"5", "1"},
			flagDefault:    []int8{1, 5},
			expectedValues: []int8{5, 1},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []int8{},
			expectedValues: []int8{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []int8{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i8s []int8
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Int8SliceVar(&i8s, "i8s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--i8s", test.input...))
			if test.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected an error; got none")
				}
				if test.expectedErr != "" && err.Error() != test.expectedErr {
					t.Fatalf("expected error to equal %q, but was: %s", test.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error; got %q", err)
			}

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			if !reflect.DeepEqual(test.expectedValues, i8s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, i8s)
			}

			int8Slice, err := f.GetInt8Slice("i8s")
			if err != nil {
				t.Fatal("got an error from GetInt8Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, int8Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, int8Slice)
			}

			int8SliceGet, err := f.Get("i8s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(int8SliceGet, int8Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, int8SliceGet)
			}
		})
	}
}
