// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	goflag "flag"
	"testing"
	"time"

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

func TestToGoflags(t *testing.T) {
	pfs := zflag.NewFlagSet("test", zflag.ContinueOnError)
	gfs := goflag.FlagSet{}
	pfs.String("StringFlag", "String value", "String flag usage")
	pfs.Int("IntFlag", 1, "Int flag usage")
	pfs.Uint("UintFlag", 2, "Uint flag usage")
	pfs.Int64("Int64Flag", 3, "Int64 flag usage")
	pfs.Uint64("Uint64Flag", 4, "Uint64 flag usage")
	pfs.Int8("Int8Flag", 5, "Int8 flag usage")
	pfs.Float64("Float64Flag", 6.0, "Float64 flag usage")
	pfs.Duration("DurationFlag", time.Second, "Duration flag usage")
	pfs.Bool("BoolFlag", true, "Bool flag usage")
	pfs.String("deprecated", "Deprecated value", "Deprecated flag usage", zflag.OptDeprecated("obsolete"))

	pfs.CopyToGoFlagSet(&gfs)

	// both flag sets share the same values
	for name, value := range map[string]string{
		"StringFlag":  "Modified String value",
		"IntFlag":     "11",
		"UintFlag":    "12",
		"Int64Flag":   "13",
		"Uint64Flag":  "14",
		"Int8Flag":    "15",
		"Float64Flag": "16.0",
		"BoolFlag":    "false",
	} {
		pf := pfs.Lookup(name)
		if pf == nil {
			t.Errorf("%s: not found in zflag flag set", name)
			continue
		}
		assertNoErr(t, pf.Value.Set(value))
	}

	pfs.VisitAll(func(pf *zflag.Flag) {
		gf := gfs.Lookup(pf.Name)
		if gf == nil {
			t.Errorf("%s: not found in Go flag set", pf.Name)
			return
		}
		assertEqual(t, pf.Value.String(), gf.Value.String())
	})

	// every Go flag must come from the zflag flag set
	gfs.VisitAll(func(gf *goflag.Flag) {
		if pfs.Lookup(gf.Name) == nil {
			t.Errorf("%s: not found in zflag flag set", gf.Name)
		}
	})

	deprecated := gfs.Lookup("deprecated")
	if deprecated == nil {
		t.Fatal("deprecated: not found in Go flag set")
	}
	assertEqual(t, "Deprecated flag usage (DEPRECATED: obsolete)", deprecated.Usage)
}
