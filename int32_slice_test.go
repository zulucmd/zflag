// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestInt32Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []int32
		input          []string
		expectedErr    string
		expectedValues []int32
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []int32{},
			expectedErr:    "",
			expectedValues: []int32{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []int32{},
			expectedErr: `invalid argument "" for "--i32s" flag: must be an integer`,
		},
		{
			name:        "invalid int32",
			input:       []string{"blabla"},
			flagDefault: []int32{},
			expectedErr: `invalid argument "blabla" for "--i32s" flag: must be an integer`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []int32{},
			expectedErr: `invalid argument "1,5" for "--i32s" flag: must be an integer`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []int32{},
			expectedValues: []int32{1, 5},
		},
		{
			name:           "with default values",
			input:          []string{"5", "1"},
			flagDefault:    []int32{1, 5},
			expectedValues: []int32{5, 1},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []int32{},
			expectedValues: []int32{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []int32{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i32s []int32
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.Int32SliceVar(&i32s, "i32s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--i32s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, i32s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, i32s)
			}

			int32Slice, err := f.GetInt32Slice("i32s")
			if err != nil {
				t.Fatal("got an error from GetInt32Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, int32Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, int32Slice)
			}

			int32SliceGet, err := f.Get("i32s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(int32SliceGet, int32Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, int32SliceGet)
			}
		})
	}
}
