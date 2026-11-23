// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestUint16Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []uint16
		input          []string
		expectedErr    string
		expectedValues []uint16
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []uint16{},
			expectedErr:    "",
			expectedValues: []uint16{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []uint16{},
			expectedErr: `invalid argument "" for "--ui16s" flag: must be a non-negative integer`,
		},
		{
			name:        "invalid uint16",
			input:       []string{"blabla"},
			flagDefault: []uint16{},
			expectedErr: `invalid argument "blabla" for "--ui16s" flag: must be a non-negative integer`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []uint16{},
			expectedErr: `invalid argument "1,5" for "--ui16s" flag: must be a non-negative integer`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []uint16{},
			expectedValues: []uint16{1, 5},
		},
		{
			name:           "overrides default values",
			input:          []string{"5", "1"},
			flagDefault:    []uint16{1, 5},
			expectedValues: []uint16{5, 1},
		},
		{
			name:           "with default values",
			input:          []string{},
			flagDefault:    []uint16{1, 5},
			expectedValues: []uint16{1, 5},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []uint16{},
			expectedValues: []uint16{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []uint16{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ui16s []uint16
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.Uint16SliceVar(&ui16s, "ui16s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ui16s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, ui16s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, ui16s)
			}

			uint16Slice, err := f.GetUint16Slice("ui16s")
			if err != nil {
				t.Fatal("got an error from GetUint16Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, uint16Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uint16Slice)
			}

			uint16SliceGet, err := f.Get("ui16s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(uint16SliceGet, uint16Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uint16SliceGet)
			}
		})
	}
}
