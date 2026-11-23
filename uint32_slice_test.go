// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestUint32Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []uint32
		input          []string
		expectedErr    string
		expectedValues []uint32
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []uint32{},
			expectedErr:    "",
			expectedValues: []uint32{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []uint32{},
			expectedErr: `invalid argument "" for "--ui32s" flag: must be a non-negative integer`,
		},
		{
			name:        "invalid uint32",
			input:       []string{"blabla"},
			flagDefault: []uint32{},
			expectedErr: `invalid argument "blabla" for "--ui32s" flag: must be a non-negative integer`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []uint32{},
			expectedErr: `invalid argument "1,5" for "--ui32s" flag: must be a non-negative integer`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []uint32{},
			expectedValues: []uint32{1, 5},
		},
		{
			name:           "overrides default values",
			input:          []string{"5", "1"},
			flagDefault:    []uint32{1, 5},
			expectedValues: []uint32{5, 1},
		},
		{
			name:           "with default values",
			input:          []string{},
			flagDefault:    []uint32{1, 5},
			expectedValues: []uint32{1, 5},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []uint32{},
			expectedValues: []uint32{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []uint32{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ui32s []uint32
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.Uint32SliceVar(&ui32s, "ui32s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ui32s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, ui32s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, ui32s)
			}

			uint32Slice, err := f.GetUint32Slice("ui32s")
			if err != nil {
				t.Fatal("got an error from GetUint32Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, uint32Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uint32Slice)
			}

			uint32SliceGet, err := f.Get("ui32s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(uint32SliceGet, uint32Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uint32SliceGet)
			}
		})
	}
}
