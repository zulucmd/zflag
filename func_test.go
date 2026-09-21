// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestFunc(t *testing.T) {
	tests := []struct {
		name           string
		input          []string
		expectedErr    string
		expectedValues []string
	}{
		{
			name:           "no value passed",
			input:          []string{},
			expectedErr:    "",
			expectedValues: []string{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			expectedErr: `invalid argument "" for "--fn" flag: func error`,
		},
		{
			name:           "no csv",
			input:          []string{"1,5"},
			expectedValues: []string{"1,5"},
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

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vals := make([]string, 0)
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.Func("fn", "usage", func(s string) error {
				if s == "" {
					return errors.New("func error")
				}

				vals = append(vals, s)

				return nil
			})
			err := f.Parse(repeatFlag("--fn", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, vals) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, vals)
			}
		})
	}
}

// Copyright 2009 The Go Authors. All rights reserved.
func TestUserDefinedFunc(t *testing.T) {
	var flags zflag.FlagSet
	flags.Init("test", zflag.ContinueOnError)
	var ss []string
	flags.Func("v", "usage", func(s string) error {
		ss = append(ss, s)
		return nil
	})
	if err := flags.Parse([]string{"-v", "1", "-v", "2", "-v=3"}); err != nil {
		t.Fatal(err)
	}
	if len(ss) != 3 {
		t.Fatal("expected 3 args; got ", len(ss))
	}
	expect := "[1 2 3]"
	if got := fmt.Sprint(ss); got != expect {
		t.Errorf("expected value %q got %q", expect, got)
	}
	// test usage
	var buf strings.Builder
	flags.SetOutput(&buf)
	_ = flags.Parse([]string{"-h"})
	if usage := buf.String(); !strings.Contains(usage, "usage") {
		t.Errorf("usage string not included: %q", usage)
	}
	// test Func error
	flags = *zflag.NewFlagSet("test", zflag.ContinueOnError)
	flags.Func("v", "usage", func(_ string) error {
		return fmt.Errorf("test error")
	})
	// flag not set, so no error
	if err := flags.Parse(nil); err != nil {
		t.Error(err)
	}
	// flag set, expect error
	if err := flags.Parse([]string{"-v", "1"}); err == nil {
		t.Error("expected error; got none")
	} else if errMsg := err.Error(); !strings.Contains(errMsg, "test error") {
		t.Errorf(`error should contain "test error"; got %q`, errMsg)
	}
}
