// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"reflect"
	"testing"
)

func repeatFlag(flag string, values ...string) (res []string) {
	res = make([]string, 0, len(values))
	for _, val := range values {
		res = append(res, fmt.Sprintf("%s=%s", flag, val))
	}

	return
}

func assertDeepEqual(t *testing.T, expected, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", expected, actual)
	}
}

func assertEqual(t *testing.T, expected, actual any) {
	t.Helper()
	assertEqualf(t, expected, actual, "expected %[1]v with type %[1]T but got %[2]v with type %[2]T", expected, actual)
}

func assertEqualf(t *testing.T, expected, actual any, msg string, fmt ...any) {
	t.Helper()
	if expected != actual {
		t.Errorf(msg, fmt...)
	}
}

func assertErr(t *testing.T, actual error) {
	t.Helper()
	if actual == nil {
		t.Fatalf("expected an error, got: %s", actual)
	}
}

func assertNoErr(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Fatalf("expected no error, got: %s", actual)
	}
}

func assertNoPanic(t *testing.T) func() {
	return func() {
		t.Helper()
		if v := recover(); v != nil {
			t.Fatal("expected no panic, got:", v)
		}
	}
}

func assertPanic(t *testing.T) func() {
	return func() {
		t.Helper()
		if v := recover(); v == nil {
			t.Fatal("expected a panic, got:", v)
		}
	}
}

func assertErrMsg(t *testing.T, expectedErrMsg string, err error) {
	t.Helper()
	assertErr(t, err)
	if err.Error() != expectedErrMsg {
		t.Fatalf("expected error to equal %q, but was: %s", expectedErrMsg, err)
	}
}

func assertNotNilf(t *testing.T, value any, msg string, fmt ...any) {
	t.Helper()
	if value == nil {
		t.Fatalf(msg, fmt...)
	}
}
