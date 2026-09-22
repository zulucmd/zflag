// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestUint64Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []uint64
		input          []string
		expectedErr    string
		expectedValues []uint64
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []uint64{},
			expectedErr:    "",
			expectedValues: []uint64{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []uint64{},
			expectedErr: `invalid argument "" for "--ui64s" flag: must be a non-negative integer`,
		},
		{
			name:        "invalid uint64",
			input:       []string{"blabla"},
			flagDefault: []uint64{},
			expectedErr: `invalid argument "blabla" for "--ui64s" flag: must be a non-negative integer`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []uint64{},
			expectedErr: `invalid argument "1,5" for "--ui64s" flag: must be a non-negative integer`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []uint64{},
			expectedValues: []uint64{1, 5},
		},
		{
			name:           "overrides default values",
			input:          []string{"5", "1"},
			flagDefault:    []uint64{1, 5},
			expectedValues: []uint64{5, 1},
		},
		{
			name:           "with default values",
			input:          []string{},
			flagDefault:    []uint64{1, 5},
			expectedValues: []uint64{1, 5},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []uint64{},
			expectedValues: []uint64{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []uint64{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ui64s []uint64
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Uint64SliceVar(&ui64s, "ui64s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ui64s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, ui64s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, ui64s)
			}

			uint64Slice, err := f.GetUint64Slice("ui64s")
			if err != nil {
				t.Fatal("got an error from GetUint64Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, uint64Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uint64Slice)
			}

			uint64SliceGet, err := f.Get("ui64s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(uint64SliceGet, uint64Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uint64SliceGet)
			}
		})
	}
}
