// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestUint8Slice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []uint8
		input          []string
		expectedErr    string
		expectedValues []uint8
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []uint8{},
			expectedErr:    "",
			expectedValues: []uint8{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []uint8{},
			expectedErr: `invalid argument "" for "--ui8s" flag: must be a non-negative integer`,
		},
		{
			name:        "invalid uint8",
			input:       []string{"blabla"},
			flagDefault: []uint8{},
			expectedErr: `invalid argument "blabla" for "--ui8s" flag: must be a non-negative integer`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []uint8{},
			expectedErr: `invalid argument "1,5" for "--ui8s" flag: must be a non-negative integer`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []uint8{},
			expectedValues: []uint8{1, 5},
		},
		{
			name:           "overrides default values",
			input:          []string{"5", "1"},
			flagDefault:    []uint8{1, 5},
			expectedValues: []uint8{5, 1},
		},
		{
			name:           "with default values",
			input:          []string{},
			flagDefault:    []uint8{1, 5},
			expectedValues: []uint8{1, 5},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []uint8{},
			expectedValues: []uint8{1, 2, 3},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []uint8{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ui8s []uint8
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Uint8SliceVar(&ui8s, "ui8s", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ui8s", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, ui8s) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, ui8s)
			}

			uint8Slice, err := f.GetUint8Slice("ui8s")
			if err != nil {
				t.Fatal("got an error from GetUint8Slice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, uint8Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uint8Slice)
			}

			uint8SliceGet, err := f.Get("ui8s")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(uint8SliceGet, uint8Slice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uint8SliceGet)
			}
		})
	}
}
