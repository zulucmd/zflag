// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	goflag "flag"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestGoflags(t *testing.T) {
	goflag.String("stringFlag", "stringFlag", "stringFlag")
	goflag.Bool("boolFlag", false, "boolFlag")
	var testxxxValue string
	goflag.StringVar(&testxxxValue, "test.xxx", "test.xxx", "it is a test flag")

	f := zflag.NewFlagSet("test", zflag.ContinueOnError)

	f.AddGoFlagSet(goflag.CommandLine)
	args := []string{"--stringFlag=bob", "--boolFlag", "-test.xxx=testvalue"}
	err := f.Parse(args)
	if err != nil {
		t.Fatal("expected no error; get", err)
	}

	getString, err := f.GetString("stringFlag")
	if err != nil {
		t.Fatal("expected no error; get", err)
	}
	if getString != "bob" {
		t.Fatalf("expected getString=bob but got getString=%s", getString)
	}
	getString2, err := f.Get("stringFlag")
	if err != nil {
		t.Fatal("expected no error; get", err)
	}
	if getString2 != "bob" {
		t.Fatalf("expected getString=bob but got getString=%s", getString2)
	}

	getBool, err := f.GetBool("boolFlag")
	if err != nil {
		t.Fatal("expected no error; get", err)
	}
	if getBool != true {
		t.Fatalf("expected getBool=true but got getBool=%v", getBool)
	}
	getBool2, err := f.Get("boolFlag")
	if err != nil {
		t.Fatal("expected no error; get", err)
	}
	if getBool2.(bool) != true {
		t.Fatalf("expected getBool2=true but got getBool2=%v", getBool2)
	}
	if !f.Parsed() {
		t.Fatal("f.Parsed() return false after f.Parse() called")
	}
	assertEqual(t, "test.xxx", testxxxValue)

	assertNoErr(t, zflag.ParseSkippedFlags(args, goflag.CommandLine))
	assertEqual(t, "testvalue", testxxxValue)

	// in fact it is useless. because `go test` called flag.Parse()
	if !goflag.CommandLine.Parsed() {
		t.Fatal("goflag.CommandLine.Parsed() return false after f.Parse() called")
	}
}
