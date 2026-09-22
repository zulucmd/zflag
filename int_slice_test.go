// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestIntSlice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []int
		input          []string
		expectedErr    string
		expectedValues []int
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []int{},
			expectedErr:    "",
			expectedValues: []int{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []int{},
			expectedErr: `invalid argument "" for "--is" flag: must be an integer`,
		},
		{
			name:        "invalid int",
			input:       []string{"blabla"},
			flagDefault: []int{},
			expectedErr: `invalid argument "blabla" for "--is" flag: must be an integer`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []int{},
			expectedErr: `invalid argument "1,5" for "--is" flag: must be an integer`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []int{},
			expectedValues: []int{1, 5},
		},
		{
			name:           "with default values",
			input:          []string{"5", "1"},
			flagDefault:    []int{1, 5},
			expectedValues: []int{5, 1},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []int{},
			expectedValues: []int{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []int{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var is []int
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.IntSliceVar(&is, "is", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--is", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, is) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, is)
			}

			intSlice, err := f.GetIntSlice("is")
			if err != nil {
				t.Fatal("got an error from GetIntSlice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, intSlice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, intSlice)
			}

			intSliceGet, err := f.Get("is")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(intSliceGet, intSlice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, intSliceGet)
			}
		})
	}
}
