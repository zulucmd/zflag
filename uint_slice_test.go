// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestUintSlice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []uint
		input          []string
		expectedErr    string
		expectedValues []uint
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []uint{},
			expectedErr:    "",
			expectedValues: []uint{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []uint{},
			expectedErr: `invalid argument "" for "--uis" flag: strconv.ParseUint: parsing "": invalid syntax`,
		},
		{
			name:        "invalid uint",
			input:       []string{"blabla"},
			flagDefault: []uint{},
			expectedErr: `invalid argument "blabla" for "--uis" flag: strconv.ParseUint: parsing "blabla": invalid syntax`,
		},
		{
			name:        "no csv",
			input:       []string{"1,5"},
			flagDefault: []uint{},
			expectedErr: `invalid argument "1,5" for "--uis" flag: strconv.ParseUint: parsing "1,5": invalid syntax`,
		},
		{
			name:           "empty defaults",
			input:          []string{"1", "5"},
			flagDefault:    []uint{},
			expectedValues: []uint{1, 5},
		},
		{
			name:           "overrides default values",
			input:          []string{"5", "1"},
			flagDefault:    []uint{1, 5},
			expectedValues: []uint{5, 1},
		},
		{
			name:           "with default values",
			input:          []string{},
			flagDefault:    []uint{1, 5},
			expectedValues: []uint{1, 5},
		},
		{
			name:           "trims input",
			input:          []string{"    1", "2    ", "   3  "},
			flagDefault:    []uint{},
			expectedValues: []uint{1, 2, 3},
		},
		{
			name:           "hex values",
			input:          []string{"0x1", "0x2", "0x10"},
			flagDefault:    []uint{},
			expectedValues: []uint{1, 2, 16},
		},
		{
			name:  "append hex value",
			input: []string{"0x1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Append("0x10")
				}
			},
			expectedValues: []uint{1, 16},
		},
		{
			name:  "replace values",
			input: []string{"5", "1"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"3"})
				}
			},
			expectedValues: []uint{3},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var uis []uint
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.UintSliceVar(&uis, "uis", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--uis", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, uis) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uis)
			}

			uintSlice, err := f.GetUintSlice("uis")
			if err != nil {
				t.Fatal("got an error from GetUintSlice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, uintSlice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uintSlice)
			}

			uintSliceGet, err := f.Get("uis")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(uintSliceGet, uintSlice) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, uintSliceGet)
			}
		})
	}
}
