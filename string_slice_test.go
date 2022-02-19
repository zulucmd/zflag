// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gowarden/zflag"
)

func TestSSValueImplementsGetter(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.StringSlice("ss", []string{"default", "values"}, "Command separated list!")
	v := f.Lookup("ss").Value

	if _, ok := v.(zflag.Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
}

func TestStringSlice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []string
		input          []string
		expectedErr    string
		expectedValues []string
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []string{},
			expectedErr:    "",
			expectedValues: []string{},
		},
		{
			name:           "empty value passed",
			input:          []string{""},
			flagDefault:    []string{},
			expectedValues: []string{""},
		},
		{
			name:           "single string",
			input:          []string{"blabla"},
			flagDefault:    []string{},
			expectedValues: []string{"blabla"},
		},
		{
			name:           "no csv",
			input:          []string{"testing,something"},
			flagDefault:    []string{},
			expectedValues: []string{"testing,something"},
		},
		{
			name:           "multiple values passed",
			input:          []string{"testing", "something", "all the strings"},
			flagDefault:    []string{},
			expectedValues: []string{"testing", "something", "all the strings"},
		},
		{
			name:           "with default values",
			input:          []string{},
			flagDefault:    []string{"testing", "0:0:0:0:0:0:0:1"},
			expectedValues: []string{"testing", "0:0:0:0:0:0:0:1"},
		},
		{
			name:           "overrides default values",
			input:          []string{"all the strings", "testing"},
			flagDefault:    []string{"testing", "0:0:0:0:0:0:0:1"},
			expectedValues: []string{"all the strings", "testing"},
		},
		{
			name:  "as slice values",
			input: []string{"testing", "all the strings"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"overridden"})
				}
			},
			expectedValues: []string{"overridden"},
		},
		{
			name:           "keeps spacing",
			input:          []string{"somestring", "        somestring", "somestring     ", "   somestring  "},
			expectedValues: []string{"somestring", "        somestring", "somestring     ", "   somestring  "},
		},
		{
			name:           "keeps new lines",
			input:          []string{"foo\nbar\nbaz\n\n\nasdasd\n\n"},
			expectedValues: []string{"foo\nbar\nbaz\n\n\nasdasd\n\n"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ss []string
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.StringSliceVar(&ss, "ss", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ss", test.input...))
			if test.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected an error; got none")
				}
				if test.expectedErr != "" && !strings.Contains(err.Error(), test.expectedErr) {
					t.Fatalf("expected error to contain %q, but was: %s", test.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error; got %q", err)
			}

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			if !reflect.DeepEqual(test.expectedValues, ss) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, ss, ss)
			}

			stringSlice, err := f.GetStringSlice("ss")
			if err != nil {
				t.Fatal("got an error from GetStringSlice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, stringSlice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, stringSlice, stringSlice)
			}

			stringSliceGet, err := f.Get("ss")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(stringSliceGet, stringSlice) {
				t.Fatalf("expected %v with type %T but got %v with type %T ", test.expectedValues, test.expectedValues, stringSliceGet, stringSliceGet)
			}
		})
	}
}
