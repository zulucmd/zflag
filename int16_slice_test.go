// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestInt16Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []int16
		input          []string
		expectedErr    string
		expectedValues []int16
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []int16{},
			expectedErr:    "",
			expectedValues: []int16{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []int16{},
			expectedErr: `invalid argument "" for "--i16s" flag: must be an integer`,
		},
		{
			name:        "invalid int16",
			input:       []string{"blabla"},
			flagDefault: []int16{},
			expectedErr: `invalid argument "blabla" for "--i16s" flag: must be an integer`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []int16{},
			expectedErr: `invalid argument "1,5" for "--i16s" flag: must be an integer`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []int16{},
			expectedValues: []int16{1, 5},
		},
		{
			name:           "with default values",
			input:          []string{"5", "1"},
			flagDefault:    []int16{1, 5},
			expectedValues: []int16{5, 1},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []int16{},
			expectedValues: []int16{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []int16{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i16s []int16
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Int16SliceVar(&i16s, "i16s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--i16s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, i16s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, i16s)
			}

			int16Slice, err := f.GetInt16Slice("i16s")
			if err != nil {
				t.Fatal("got an error from GetInt16Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, int16Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, int16Slice)
			}

			int16SliceGet, err := f.Get("i16s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(int16SliceGet, int16Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, int16SliceGet)
			}

			expectedString := fmt.Sprintf("%v", test.expectedValues)
			if got := f.Lookup("i16s").Value.String(); got != expectedString {
				t.Fatalf("expected String() to be %q, but was %q", expectedString, got)
			}
		})
	}
}
