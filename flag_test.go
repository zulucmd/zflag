// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/zulucmd/zflag/v2"
)

var (
	_                            = zflag.Bool("test_bool", false, "bool value")
	_                            = zflag.Int("test_int", 0, "int value")
	_                            = zflag.Int64("test_int64", 0, "int64 value")
	_                            = zflag.Uint("test_uint", 0, "uint value")
	_                            = zflag.Uint64("test_uint64", 0, "uint64 value")
	_                            = zflag.String("test_string", "0", "string value")
	_                            = zflag.Float64("test_float64", 0, "float64 value")
	_                            = zflag.Duration("test_duration", 0, "time.Duration value")
	_                            = zflag.Int("test_optional_int", 0, "optional int value")
	normalizeFlagNameInvocations = 0
)

func TestCmdVars(_ *testing.T) {
	var tbool bool
	zflag.BoolVar(&tbool, "bool_var", false, "bool value")

	var tint int
	zflag.IntVar(&tint, "int_var", 0, "int value")

	var tint64 int64
	zflag.Int64Var(&tint64, "int64_var", 0, "int64 value")

	var tUint uint
	zflag.UintVar(&tUint, "uint_var", 0, "uint value")

	var tUint64 uint64
	zflag.Uint64Var(&tUint64, "uint64_var", 0, "uint64 value")

	var tString string
	zflag.StringVar(&tString, "string_var", "0", "string value")

	var tFloat64 float64
	zflag.Float64Var(&tFloat64, "float64_var", 0, "float64 value")

	var tDuration time.Duration
	zflag.DurationVar(&tDuration, "duration_var", 0, "time.Duration value")

	var tOptInt int
	zflag.IntVar(&tOptInt, "optional_int_var", 0, "optional int value")

	var tBoolSlice []bool
	_ = zflag.BoolSlice("bool_slice", []bool{}, "usage")
	zflag.BoolSliceVar(&tBoolSlice, "bool_slice_var", []bool{}, "usage")

	var tBytesHex []byte
	_ = zflag.BytesHex("bytes_hex", nil, "usage")
	zflag.BytesHexVar(&tBytesHex, "bytes_hex_var", nil, "usage")
	var tBytesBase64 []byte
	_ = zflag.BytesBase64("bytes_base64", nil, "usage")
	zflag.BytesBase64Var(&tBytesBase64, "bytes_base64_var", nil, "usage")

	var tc128 complex128
	_ = zflag.Complex128("c128", complex(0, 0), "usage")
	zflag.Complex128Var(&tc128, "c128_var", complex(0, 0), "usage")

	var tc128s []complex128
	_ = zflag.Complex128Slice("c128s", []complex128{}, "usage")
	zflag.Complex128SliceVar(&tc128s, "c128s_var", []complex128{}, "usage")

	zflag.Func("func", "", func(_ string) error {
		return nil
	})

	var tCount int
	_ = zflag.Count("count", "")
	zflag.CountVar(&tCount, "count_var", "")

	var tDurations []time.Duration
	_ = zflag.DurationSlice("durations", []time.Duration{}, "usage")
	zflag.DurationSliceVar(&tDurations, "durations_var", []time.Duration{}, "usage")

	var tTime time.Time
	_ = zflag.Time("time", time.Time{}, []string{time.RFC3339}, "usage")
	zflag.TimeVar(&tTime, "time_var", time.Time{}, []string{time.RFC3339}, "usage")
}

func boolString(s string) string {
	if s == "0" {
		return "false"
	}
	return "true"
}

func TestEverything(t *testing.T) {
	m := make(map[string]*zflag.Flag)
	desired := "0"
	visitor := func(f *zflag.Flag) {
		if len(f.Name) > 5 && f.Name[0:5] == "test_" {
			m[f.Name] = f
			ok := false
			switch {
			case f.Value.String() == desired:
				ok = true
			case f.Name == "test_bool" && f.Value.String() == boolString(desired):
				ok = true
			case f.Name == "test_duration" && f.Value.String() == desired+"s":
				ok = true
			}
			if !ok {
				t.Error("Visit: bad value", f.Value.String(), "for", f.Name)
			}
		}
	}
	assertEqual(t, true, zflag.CommandLine.HasFlags())
	zflag.VisitAll(visitor)
	if len(m) != 9 {
		t.Error("VisitAll misses some flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	m = make(map[string]*zflag.Flag)
	zflag.Visit(visitor)
	if len(m) != 0 {
		t.Errorf("Visit sees unset flags")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now set all flags
	_ = zflag.Set("test_bool", "true")
	_ = zflag.Set("test_int", "1")
	_ = zflag.Set("test_int64", "1")
	_ = zflag.Set("test_uint", "1")
	_ = zflag.Set("test_uint64", "1")
	_ = zflag.Set("test_string", "1")
	_ = zflag.Set("test_float64", "1")
	_ = zflag.Set("test_duration", "1s")
	_ = zflag.Set("test_optional_int", "1")
	desired = "1"
	zflag.Visit(visitor)
	if len(m) != 9 {
		t.Error("Visit fails after set")
		for k, v := range m {
			t.Log(k, *v)
		}
	}
	// Now test they're visited in sort order.
	var flagNames []string
	zflag.Visit(func(f *zflag.Flag) { flagNames = append(flagNames, f.Name) })
	if !sort.StringsAreSorted(flagNames) {
		t.Errorf("flag names not sorted: %v", flagNames)
	}
}

func TestUsage(t *testing.T) {
	called := false
	zflag.ResetForTesting(func() { called = true })
	if zflag.CommandLine.Parse([]string{"--x"}) == nil {
		t.Error("parse did not fail for unknown flag")
	}
	if !called {
		t.Error("did not call Usage for unknown flag")
	}
}

func TestAddFlagSet(t *testing.T) {
	oldSet := zflag.NewFlagSet("old", zflag.ContinueOnError)
	newSet := zflag.NewFlagSet("new", zflag.ContinueOnError)

	oldSet.String("flag1", "flag1", "flag1")
	oldSet.String("flag2", "flag2", "flag2")

	newSet.String("flag2", "flag2", "flag2")
	newSet.String("flag3", "flag3", "flag3")

	oldSet.AddFlagSet(newSet)

	if len(zflag.GetFlagFormalField(oldSet)) != 3 {
		t.Errorf("Unexpected result adding a FlagSet to a FlagSet %v", oldSet)
	}
}

func TestAddFlag(t *testing.T) {
	fs := zflag.NewFlagSet("adding-flags", zflag.ContinueOnError)
	fs.AddFlag(&zflag.Flag{
		Name:      "a-flag",
		Shorthand: 'a',
		Usage:     "the usage",
	})

	if len(zflag.GetFlagFormalField(fs)) != 1 {
		t.Errorf("Unexpected result adding a Flag to a FlagSet %v", fs)
	}
}

func TestRemoveFlag(t *testing.T) {
	fs := zflag.NewFlagSet("removing-flags", zflag.ContinueOnError)
	fs.String("string1", "a", "enter a string1")
	fs.String("string2", "b", "enter a string2")

	fs.RemoveFlag("string1")

	formal := zflag.GetFlagFormalField(fs)
	if len(formal) != 1 {
		t.Errorf("Flagset %v should only have 1 formal flag now, but has %d.", fs, len(formal))
	}
	fs.VisitAll(func(f *zflag.Flag) {
		if f.Name == "string1" {
			t.Errorf("Flag string1 was not removed from %v formal flags.", fs)
		}
	})
}

func TestAnnotation(t *testing.T) {
	f := zflag.NewFlagSet("shorthand", zflag.ContinueOnError)

	f.String("stringa", "", "string value", zflag.OptShorthand('a'), zflag.OptAnnotation("key", nil))
	if annotation := f.Lookup("stringa").Annotations["key"]; annotation != nil {
		t.Errorf("Unexpected annotation: %v", annotation)
	}

	f.String("stringb", "", "string2 value", zflag.OptShorthand('b'), zflag.OptAnnotation("key", []string{"value1"}))
	stringb := f.Lookup("stringb")
	if annotation := stringb.Annotations["key"]; !reflect.DeepEqual(annotation, []string{"value1"}) {
		t.Errorf("Unexpected annotation: %v", annotation)
	}

	stringb.SetAnnotation("key", []string{"value2"})
	if annotation := stringb.Annotations["key"]; !reflect.DeepEqual(annotation, []string{"value2"}) {
		t.Errorf("Unexpected annotation: %v", annotation)
	}
}

func TestName(t *testing.T) {
	flagSetName := "bob"
	f := zflag.NewFlagSet(flagSetName, zflag.ContinueOnError)

	givenName := f.Name()
	if givenName != flagSetName {
		t.Errorf("Unexpected result when retrieving a FlagSet's name: expected %s, but found %s", flagSetName, givenName)
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name                      string
		args                      []string
		expectedError             string
		IgnoreRequiredFlagsErrors bool
	}{
		{
			name:                      "errors when required flag not there",
			args:                      []string{"--string=hello", "some-arg"},
			expectedError:             `required flag(s) "--required-int", "--required-string" not set`,
			IgnoreRequiredFlagsErrors: false,
		},
		{
			name:                      "errors when required flag not there",
			args:                      []string{"--required-string=hello", "some-arg"},
			expectedError:             `required flag(s) "--required-int" not set`,
			IgnoreRequiredFlagsErrors: false,
		},
		{
			name:                      "does not error when required flags are there",
			args:                      []string{"--required-int=4", "--required-string=hello", "some-arg"},
			expectedError:             "",
			IgnoreRequiredFlagsErrors: false,
		},
		{
			name:                      "does not error when required flags are not there and required flags are not ignored",
			args:                      []string{"some-arg"},
			expectedError:             "",
			IgnoreRequiredFlagsErrors: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			if f.Parsed() {
				t.Error("f.Parse() = true before Parse")
			}
			_ = f.String("string", "0", "string value")
			_ = f.String("required-string", "0", "required string value", zflag.OptRequired())
			_ = f.Int("required-int", 0, "required int value", zflag.OptRequired())
			f.ParseErrorsAllowList.RequiredFlags = tt.IgnoreRequiredFlagsErrors

			err := f.Parse(tt.args)
			if err == nil && tt.expectedError != "" {
				t.Errorf("Parse() got no error but expected: %v", tt.expectedError)
				return
			}

			if err != nil && tt.expectedError == "" {
				t.Errorf("Parse() expected no error, got: %v", err)
				return
			}

			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("Parse() error = %v, wantErr: %v", err, tt.expectedError)
			}
		})
	}
}

func testParse(f *zflag.FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	boolFlag := f.Bool("bool", false, "bool value")
	bool2Flag := f.Bool("bool2", false, "bool2 value")
	bool3Flag := f.Bool("bool3", false, "bool3 value")
	intFlag := f.Int("int", 0, "int value")
	int8Flag := f.Int8("int8", 0, "int value")
	int16Flag := f.Int16("int16", 0, "int value")
	int32Flag := f.Int32("int32", 0, "int value")
	int64Flag := f.Int64("int64", 0, "int64 value")
	uintFlag := f.Uint("uint", 0, "uint value")
	uint8Flag := f.Uint8("uint8", 0, "uint value")
	uint16Flag := f.Uint16("uint16", 0, "uint value")
	uint32Flag := f.Uint32("uint32", 0, "uint value")
	uint64Flag := f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	float32Flag := f.Float32("float32", 0, "float32 value")
	float64Flag := f.Float64("float64", 0, "float64 value")
	ipFlag := f.IP("ip", net.ParseIP("127.0.0.1"), "ip value")
	maskFlag := f.IPMask("mask", zflag.ParseIPv4Mask("0.0.0.0"), "mask value")
	durationFlag := f.Duration("duration", 5*time.Second, "time.Duration value")
	extra := "one-extra-argument"
	args := []string{
		"--bool",
		"--bool2=true",
		"--bool3=false",
		"--int=22",
		"--int8=-8",
		"--int16=-16",
		"--int32=-32",
		"--int64=0x23",
		"--uint", "24",
		"--uint8=8",
		"--uint16=16",
		"--uint32=32",
		"--uint64=25",
		"--string=hello",
		"--float32=-172e12",
		"--float64=2718e28",
		"--ip=10.11.12.13",
		"--mask=255.255.255.0",
		"--duration=2m",
		extra,
	}
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if *boolFlag != true {
		t.Error("bool flag should be true, is ", *boolFlag)
	}
	if v, err := f.GetBool("bool"); err != nil || v != *boolFlag {
		t.Error("GetBool does not work.")
	}
	if v, err := f.Get("bool"); err != nil || v.(bool) != *boolFlag {
		t.Error("GetBool does not work.")
	}
	if *bool2Flag != true {
		t.Error("bool2 flag should be true, is ", *bool2Flag)
	}
	if *bool3Flag != false {
		t.Error("bool3 flag should be false, is ", *bool2Flag)
	}
	if *intFlag != 22 {
		t.Error("int flag should be 22, is ", *intFlag)
	}
	if v, err := f.GetInt("int"); err != nil || v != *intFlag {
		t.Error("GetInt does not work.")
	}
	if v, err := f.Get("int"); err != nil || v.(int) != *intFlag {
		t.Error("Get does not work.")
	}
	if *int8Flag != -8 {
		t.Error("int8 flag should be 0x23, is ", *int8Flag)
	}
	if *int16Flag != -16 {
		t.Error("int16 flag should be -16, is ", *int16Flag)
	}
	if v, err := f.GetInt8("int8"); err != nil || v != *int8Flag {
		t.Error("GetInt8 does not work.")
	}
	if v, err := f.Get("int8"); err != nil || v.(int8) != *int8Flag {
		t.Error("Get does not work.")
	}
	if v, err := f.GetInt16("int16"); err != nil || v != *int16Flag {
		t.Error("GetInt16 does not work.")
	}
	if v, err := f.Get("int16"); err != nil || v.(int16) != *int16Flag {
		t.Error("Get does not work.")
	}
	if *int32Flag != -32 {
		t.Error("int32 flag should be 0x23, is ", *int32Flag)
	}
	if v, err := f.GetInt32("int32"); err != nil || v != *int32Flag {
		t.Error("GetInt32 does not work.")
	}
	if v, err := f.Get("int32"); err != nil || v != *int32Flag {
		t.Error("Get does not work.")
	}
	if *int64Flag != 0x23 {
		t.Error("int64 flag should be 0x23, is ", *int64Flag)
	}
	if v, err := f.GetInt64("int64"); err != nil || v != *int64Flag {
		t.Error("GetInt64 does not work.")
	}
	if v, err := f.Get("int64"); err != nil || v != *int64Flag {
		t.Error("Get does not work.")
	}
	if *uintFlag != 24 {
		t.Error("uint flag should be 24, is ", *uintFlag)
	}
	if v, err := f.GetUint("uint"); err != nil || v != *uintFlag {
		t.Error("GetUint does not work.")
	}
	if v, err := f.Get("uint"); err != nil || v.(uint) != *uintFlag {
		t.Error("Get does not work.")
	}
	if *uint8Flag != 8 {
		t.Error("uint8 flag should be 8, is ", *uint8Flag)
	}
	if v, err := f.GetUint8("uint8"); err != nil || v != *uint8Flag {
		t.Error("GetUint8 does not work.")
	}
	if v, err := f.Get("uint8"); err != nil || v.(uint8) != *uint8Flag {
		t.Error("Get does not work.")
	}
	if *uint16Flag != 16 {
		t.Error("uint16 flag should be 16, is ", *uint16Flag)
	}
	if v, err := f.GetUint16("uint16"); err != nil || v != *uint16Flag {
		t.Error("GetUint16 does not work.")
	}
	if v, err := f.Get("uint16"); err != nil || v.(uint16) != *uint16Flag {
		t.Error("Get does not work.")
	}
	if *uint32Flag != 32 {
		t.Error("uint32 flag should be 32, is ", *uint32Flag)
	}
	if v, err := f.GetUint32("uint32"); err != nil || v != *uint32Flag {
		t.Error("GetUint32 does not work.")
	}
	if v, err := f.Get("uint32"); err != nil || v.(uint32) != *uint32Flag {
		t.Error("Get does not work.")
	}
	if *uint64Flag != 25 {
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	}
	if v, err := f.GetUint64("uint64"); err != nil || v != *uint64Flag {
		t.Error("GetUint64 does not work.")
	}
	if v, err := f.Get("uint64"); err != nil || v.(uint64) != *uint64Flag {
		t.Error("Get does not work.")
	}
	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is ", *stringFlag)
	}
	if v, err := f.GetString("string"); err != nil || v != *stringFlag {
		t.Error("GetString does not work.")
	}
	if v, err := f.Get("string"); err != nil || v.(string) != *stringFlag {
		t.Error("Get does not work.")
	}
	if *float32Flag != -172e12 {
		t.Error("float32 flag should be -172e12, is ", *float32Flag)
	}
	if v, err := f.GetFloat32("float32"); err != nil || v != *float32Flag {
		t.Errorf("GetFloat32 returned %v but float32Flag was %v", v, *float32Flag)
	}
	if v, err := f.Get("float32"); err != nil || v.(float32) != *float32Flag {
		t.Errorf("Get returned %v but float32Flag was %v", v, *float32Flag)
	}
	if *float64Flag != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	}
	if v, err := f.GetFloat64("float64"); err != nil || v != *float64Flag {
		t.Errorf("GetFloat64 returned %v but float64Flag was %v", v, *float64Flag)
	}
	if v, err := f.Get("float64"); err != nil || v.(float64) != *float64Flag {
		t.Errorf("Get returned %v but float64Flag was %v", v, *float64Flag)
	}
	if !ipFlag.Equal(net.ParseIP("10.11.12.13")) {
		t.Error("ip flag should be 10.11.12.13, is ", *ipFlag)
	}
	if v, err := f.GetIP("ip"); err != nil || !v.Equal(*ipFlag) {
		t.Errorf("GetIP returned %v but ipFlag was %v", v, *ipFlag)
	}
	if v, err := f.Get("ip"); err != nil || !v.(net.IP).Equal(*ipFlag) {
		t.Errorf("GetIP returned %v but ipFlag was %v", v, *ipFlag)
	}
	if maskFlag.String() != zflag.ParseIPv4Mask("255.255.255.0").String() {
		t.Error("mask flag should be 255.255.255.0, is ", maskFlag.String())
	}
	if v, err := f.GetIPMask("mask"); err != nil || v.String() != maskFlag.String() {
		t.Errorf("GetIP returned %v maskFlag was %v error was %v", v, *maskFlag, err)
	}
	if v, err := f.Get("mask"); err != nil || v.(net.IPMask).String() != maskFlag.String() {
		t.Errorf("Get returned %v maskFlag was %v error was %v", v, *maskFlag, err)
	}
	if *durationFlag != 2*time.Minute {
		t.Error("duration flag should be 2m, is ", *durationFlag)
	}
	if v, err := f.GetDuration("duration"); err != nil || v != *durationFlag {
		t.Error("GetDuration does not work.")
	}
	if v, err := f.Get("duration"); err != nil || v.(time.Duration) != *durationFlag {
		t.Error("Get does not work.")
	}
	if _, err := f.GetInt("duration"); err == nil {
		t.Error("GetInt parsed a time.Duration?!?!")
	}
	if len(f.Args()) != 1 {
		t.Error("expected one argument, got", len(f.Args()))
	} else if f.Args()[0] != extra {
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	}
}

func testParseAll(f *zflag.FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	f.Bool("boola", false, "bool value", zflag.OptShorthand('a'))
	f.Bool("boolb", false, "bool2 value", zflag.OptShorthand('b'))
	f.Bool("boolc", false, "bool3 value", zflag.OptShorthand('c'))
	f.Bool("boold", false, "bool4 value", zflag.OptShorthand('d'))
	f.String("stringa", "0", "string value", zflag.OptShorthand('s'))
	f.String("stringz", "0", "string value", zflag.OptShorthand('z'))
	f.String("stringy", "0", "string value", zflag.OptShorthand('y'))
	args := []string{
		"-ab",
		"-cs=xx",
		"--stringz=something",
		"-d=true",
		"-y",
		"ee",
	}
	want := []string{
		"boola",
		"boolb",
		"boolc",
		"stringa", "xx",
		"stringz", "something",
		"boold", "true",
		"stringy", "ee",
	}
	got := []string{}
	store := func(flag *zflag.Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.ParseAll() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}
}

func testParseWithUnknownFlags(f *zflag.FlagSet, t *testing.T) {
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	f.ParseErrorsAllowList.UnknownFlags = true

	f.Bool("boola", false, "bool value", zflag.OptShorthand('a'))
	f.Bool("boolb", false, "bool2 value", zflag.OptShorthand('b'))
	f.Bool("boolc", false, "bool3 value", zflag.OptShorthand('c'))
	f.Bool("boold", false, "bool4 value", zflag.OptShorthand('d'))
	f.Bool("boole", false, "bool4 value", zflag.OptShorthand('e'))
	f.String("stringa", "0", "string value", zflag.OptShorthand('s'))
	f.String("stringz", "0", "string value", zflag.OptShorthand('z'))
	f.String("stringy", "0", "string value", zflag.OptShorthand('y'))
	f.String("stringo", "0", "string value", zflag.OptShorthand('o'))
	args := []string{
		"-ab",
		"-cs=xx",
		"--stringz=something",
		"--unknown1",
		"unknown1Value",
		"-d=true",
		"--unknown2=unknown2Value",
		"-u=unknown3Value",
		"-p",
		"unknown4Value",
		"-q", // another unknown with bool value
		"-y",
		"ee",
		"--unknown7=unknown7value",
		"--stringo=ovalue",
		"--unknown8=unknown8value",
		"--boole",
		"--unknown6",
		"",
		"-uuuuu",
		"",
		"--unknown10",
		"--unknown11",
	}
	want := []string{
		"boola",
		"boolb",
		"boolc",
		"stringa", "xx",
		"stringz", "something",
		"boold", "true",
		"stringy", "ee",
		"stringo", "ovalue",
		"boole", "true",
	}
	wantUnknowns := []string{
		"--unknown1", "unknown1Value",
		"--unknown2=unknown2Value",
		"-u=unknown3Value",
		"-p", "unknown4Value",
		"-q",
		"--unknown7=unknown7value",
		"--unknown8=unknown8value",
		"--unknown6", "",
		"-uuuuu",
		"--unknown10",
		"--unknown11",
	}
	got := []string{}
	store := func(flag *zflag.Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parsed() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.Parse() failed to parse with unknown flags")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}
	gotUnknowns := f.GetUnknownFlags()
	if !reflect.DeepEqual(gotUnknowns, wantUnknowns) {
		t.Errorf("f.Parse() failed to enumerate the unknown args args")
		t.Errorf("Got:  %v", gotUnknowns)
		t.Errorf("Want: %v", wantUnknowns)
	}
}

func TestShorthand(t *testing.T) {
	f := zflag.NewFlagSet("shorthand", zflag.ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	boolaFlag := f.Bool("boola", false, "bool value", zflag.OptShorthand('a'))
	boolbFlag := f.Bool("boolb", false, "bool2 value", zflag.OptShorthand('b'))
	boolcFlag := f.Bool("boolc", false, "bool3 value", zflag.OptShorthand('c'))
	booldFlag := f.Bool("boold", false, "bool4 value", zflag.OptShorthand('d'))
	booleFlag := f.Bool("boole", false, "bool5 value", zflag.OptShorthand('e'))
	stringaFlag := f.String("stringa", "0", "string value", zflag.OptShorthand('s'))
	stringzFlag := f.String("stringz", "0", "string value", zflag.OptShorthand('z'))
	extra := "interspersed-argument"
	notaflag := "--i-look-like-a-flag"
	args := []string{
		"-abe",
		extra,
		"-cs",
		"hello",
		"-z=something",
		"-d=true",
		"--",
		notaflag,
	}
	f.SetOutput(ioutil.Discard)
	if err := f.Parse(args); err != nil {
		t.Error("expected no error, got", err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if *boolaFlag != true {
		t.Error("boola flag should be true, is ", *boolaFlag)
	}
	if *boolbFlag != true {
		t.Error("boolb flag should be true, is ", *boolbFlag)
	}
	if *boolcFlag != true {
		t.Error("boolc flag should be true, is ", *boolcFlag)
	}
	if *booldFlag != true {
		t.Error("boold flag should be true, is ", *booldFlag)
	}
	if *booleFlag != true {
		t.Error("boole flag should be true, is ", *booleFlag)
	}
	if *stringaFlag != "hello" {
		t.Error("stringa flag should be `hello`, is ", *stringaFlag)
	}
	if *stringzFlag != "something" {
		t.Error("stringz flag should be `something`, is ", *stringzFlag)
	}
	switch {
	case len(f.Args()) != 2:
		t.Error("expected one argument, got", len(f.Args()))
	case f.Args()[0] != extra:
		t.Errorf("expected argument %q got %q", extra, f.Args()[0])
	case f.Args()[1] != notaflag:
		t.Errorf("expected argument %q got %q", notaflag, f.Args()[1])
	}
	if f.ArgsLenAtDash() != 1 {
		t.Errorf("expected argsLenAtDash %d got %d", f.ArgsLenAtDash(), 1)
	}
}

func TestShorthandOnly(t *testing.T) {
	f := zflag.NewFlagSet("shorthand", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	boolFlag := f.Bool("bool", false, "bool value", zflag.OptShorthandOnly(), zflag.OptShorthand('1'))
	args := []string{
		"--bool",
	}
	if err := f.Parse(args); err != nil {
		t.Error("expected no error, got ", err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if *boolFlag {
		t.Error("ShorthandOnly boolFlag should be false when passed in long form")
	}
}

func TestShorthandLookup(t *testing.T) {
	f := zflag.NewFlagSet("shorthand", zflag.ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	f.Bool("boola", false, "bool value", zflag.OptShorthand('a'))
	f.Bool("boolb", false, "bool2 value", zflag.OptShorthand('b'))
	f.Bool("boolö", false, "bool2 value", zflag.OptShorthand('ö'))
	args := []string{
		"-ab",
	}
	f.SetOutput(ioutil.Discard)
	err := f.Parse(args)
	assertNoErr(t, err)
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	flag := f.ShorthandLookup('a')
	assertNotNilf(t, flag, "f.ShorthandLookup('a') returned nil")
	assertEqualf(t, flag.Name, "boola", `f.ShorthandLookup('a') found %q instead of "boola"`, flag.Name)

	flag = f.ShorthandLookup('d')
	assertNotNilf(t, flag, "f.ShorthandLookup('d') did not return nil")
	flag = f.ShorthandLookup('ö')
	assertEqualf(t, flag.Name, "boolö", `f.ShorthandLookup('ö') found %q instead of "boolö"`, flag.Name)

	flag = f.ShorthandLookupStr("a")
	assertNotNilf(t, flag, `f.ShorthandLookupStrStr("a") returned nil`)
	assertEqualf(t, flag.Name, "boola", `f.ShorthandLookupStr('a') found %q instead of "boola"`, flag.Name)

	flag = f.ShorthandLookupStr("d")
	assertNotNilf(t, flag, `f.ShorthandLookupStr("d") did not return nil`)
	flag = f.ShorthandLookupStr("ö")
	assertEqualf(t, flag.Name, "boolö", `f.ShorthandLookupStr("ö") found %q instead of "boolö"`, flag.Name)

	func() {
		defer assertPanic(t)()
		flag = f.ShorthandLookupStr("aa")
	}()

	func() {
		defer assertNoPanic(t)()
		flag = f.ShorthandLookupStr("")
	}()
}

func TestParse(t *testing.T) {
	zflag.ResetForTesting(func() { t.Error("bad parse") })
	testParse(zflag.CommandLine, t)
}

func TestParseAll(t *testing.T) {
	zflag.ResetForTesting(func() { t.Error("bad parse") })
	testParseAll(zflag.CommandLine, t)
}

func TestIgnoreUnknownFlags(t *testing.T) {
	zflag.ResetForTesting(func() { t.Error("bad parse") })
	testParseWithUnknownFlags(zflag.CommandLine, t)
}

func TestFlagSetParse(t *testing.T) {
	testParse(zflag.NewFlagSet("test", zflag.ContinueOnError), t)
}

func TestChangedHelper(t *testing.T) {
	f := zflag.NewFlagSet("changedtest", zflag.ContinueOnError)
	f.Bool("changed", false, "changed bool")
	f.Bool("settrue", true, "true to true")
	f.Bool("setfalse", false, "false to false")
	f.Bool("unchanged", false, "unchanged bool")

	args := []string{"--changed", "--settrue", "--setfalse=false"}
	if err := f.Parse(args); err != nil {
		t.Error("f.Parse() = false after Parse")
	}
	if !f.Changed("changed") {
		t.Errorf("--changed wasn't changed!")
	}
	if !f.Changed("settrue") {
		t.Errorf("--settrue wasn't changed!")
	}
	if !f.Changed("setfalse") {
		t.Errorf("--setfalse wasn't changed!")
	}
	if f.Changed("unchanged") {
		t.Errorf("--unchanged was changed!")
	}
	if f.Changed("invalid") {
		t.Errorf("--invalid was changed!")
	}
	if f.ArgsLenAtDash() != -1 {
		t.Errorf("Expected argsLenAtDash: %d but got %d", -1, f.ArgsLenAtDash())
	}
}

func replaceSeparators(name string, from []string, to string) string {
	result := name
	for _, sep := range from {
		result = strings.ReplaceAll(result, sep, to)
	}
	// Type convert to indicate normalization has been done.
	return result
}

func wordSepNormalizeFunc(_ *zflag.FlagSet, name string) zflag.NormalizedName {
	seps := []string{"-", "_"}
	name = replaceSeparators(name, seps, ".")
	normalizeFlagNameInvocations++

	return zflag.NormalizedName(name)
}

func testWordSepNormalizedNames(args []string, t *testing.T) {
	f := zflag.NewFlagSet("normalized", zflag.ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	withDashFlag := f.Bool("with-dash-flag", false, "bool value")
	// Set this after some flags have been added and before others.
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	withUnderFlag := f.Bool("with_under_flag", false, "bool value")
	withBothFlag := f.Bool("with-both_flag", false, "bool value")
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if *withDashFlag != true {
		t.Error("withDashFlag flag should be true, is ", *withDashFlag)
	}
	if *withUnderFlag != true {
		t.Error("withUnderFlag flag should be true, is ", *withUnderFlag)
	}
	if *withBothFlag != true {
		t.Error("withBothFlag flag should be true, is ", *withBothFlag)
	}
}

func TestWordSepNormalizedNames(t *testing.T) {
	args := []string{
		"--with-dash-flag",
		"--with-under-flag",
		"--with-both-flag",
	}
	testWordSepNormalizedNames(args, t)

	args = []string{
		"--with_dash_flag",
		"--with_under_flag",
		"--with_both_flag",
	}
	testWordSepNormalizedNames(args, t)

	args = []string{
		"--with-dash_flag",
		"--with-under_flag",
		"--with-both_flag",
	}
	testWordSepNormalizedNames(args, t)
}

func aliasAndWordSepFlagNames(_ *zflag.FlagSet, name string) zflag.NormalizedName {
	seps := []string{"-", "_"}

	oldName := replaceSeparators("old-valid_flag", seps, ".")
	newName := replaceSeparators("valid-flag", seps, ".")

	name = replaceSeparators(name, seps, ".")
	if name == oldName {
		name = newName
	}

	return zflag.NormalizedName(name)
}

func TestCustomNormalizedNames(t *testing.T) {
	f := zflag.NewFlagSet("normalized", zflag.ContinueOnError)
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}

	validFlag := f.Bool("valid-flag", false, "bool value")
	f.SetNormalizeFunc(aliasAndWordSepFlagNames)
	someOtherFlag := f.Bool("some-other-flag", false, "bool value")

	args := []string{"--old_valid_flag", "--some-other_flag"}
	if err := f.Parse(args); err != nil {
		t.Fatal(err)
	}

	if *validFlag != true {
		t.Errorf("validFlag is %v even though we set the alias --old_valid_falg", *validFlag)
	}
	if *someOtherFlag != true {
		t.Error("someOtherFlag should be true, is ", *someOtherFlag)
	}
}

// Every flag we add, the name (displayed also in usage) should normalized
func TestNormalizationFuncShouldChangeFlagName(t *testing.T) {
	// Test normalization after addition
	f := zflag.NewFlagSet("normalized", zflag.ContinueOnError)

	f.Bool("valid_flag", false, "bool value")
	if f.Lookup("valid_flag").Name != "valid_flag" {
		t.Error("The new flag should have the name 'valid_flag' instead of ", f.Lookup("valid_flag").Name)
	}

	f.SetNormalizeFunc(wordSepNormalizeFunc)
	if f.Lookup("valid_flag").Name != "valid.flag" {
		t.Error("The new flag should have the name 'valid.flag' instead of ", f.Lookup("valid_flag").Name)
	}

	// Test normalization before addition
	f = zflag.NewFlagSet("normalized", zflag.ContinueOnError)
	f.SetNormalizeFunc(wordSepNormalizeFunc)

	f.Bool("valid_flag", false, "bool value")
	if f.Lookup("valid_flag").Name != "valid.flag" {
		t.Error("The new flag should have the name 'valid.flag' instead of ", f.Lookup("valid_flag").Name)
	}
}

// Related to https://github.com/spf13/cobra/issues/521.
func TestNormalizationSharedFlags(t *testing.T) {
	f := zflag.NewFlagSet("set f", zflag.ContinueOnError)
	g := zflag.NewFlagSet("set g", zflag.ContinueOnError)
	nfunc := wordSepNormalizeFunc
	testName := "valid_flag"
	normName := nfunc(nil, testName)
	if testName == string(normName) {
		t.Error("TestNormalizationSharedFlags meaningless: the original and normalized flag names are identical:", testName)
	}

	f.Bool(testName, false, "bool value")
	g.AddFlagSet(f)

	f.SetNormalizeFunc(nfunc)
	g.SetNormalizeFunc(nfunc)

	if len(zflag.GetFlagFormalField(f)) != 1 {
		t.Error("Normalizing flags should not result in duplications in the flag set:", zflag.GetFlagFormalField(f))
	}
	if zflag.GetFlagOrderedFormalField(f)[0].Name != string(normName) {
		t.Error("Flag name not normalized")
	}
	for k := range zflag.GetFlagFormalField(f) {
		if k != "valid.flag" {
			t.Errorf("The key in the flag map should have been normalized: wanted \"%s\", got \"%s\" instead", normName, k)
		}
	}

	if !reflect.DeepEqual(zflag.GetFlagFormalField(f), zflag.GetFlagFormalField(g)) || !reflect.DeepEqual(zflag.GetFlagOrderedFormalField(f), zflag.GetFlagOrderedFormalField(g)) {
		t.Error("Two flag sets sharing the same flags should stay consistent after being normalized. Original set:", zflag.GetFlagFormalField(f), "Duplicate set:", zflag.GetFlagFormalField(g))
	}
}

func TestNormalizationSetFlags(t *testing.T) {
	f := zflag.NewFlagSet("normalized", zflag.ContinueOnError)
	nfunc := wordSepNormalizeFunc
	testName := "valid_flag"
	normName := nfunc(nil, testName)
	if testName == string(normName) {
		t.Error("TestNormalizationSetFlags meaningless: the original and normalized flag names are identical:", testName)
	}

	f.Bool(testName, false, "bool value")
	_ = f.Set(testName, "true")
	f.SetNormalizeFunc(nfunc)

	if len(zflag.GetFlagFormalField(f)) != 1 {
		t.Error("Normalizing flags should not result in duplications in the flag set:", zflag.GetFlagFormalField(f))
	}
	if zflag.GetFlagOrderedFormalField(f)[0].Name != string(normName) {
		t.Error("Flag name not normalized")
	}
	for k := range zflag.GetFlagFormalField(f) {
		if k != "valid.flag" {
			t.Errorf("The key in the flag map should have been normalized: wanted \"%s\", got \"%s\" instead", normName, k)
		}
	}

	if !reflect.DeepEqual(zflag.GetFlagFormalField(f), zflag.GetActual(f)) {
		t.Error("The map of set flags should get normalized. Formal:", zflag.GetFlagFormalField(f), "Actual:", zflag.GetActual(f))
	}
}

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *flagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *flagVar) Type() string {
	return "flagVar"
}

func TestGroups(t *testing.T) {
	var fs zflag.FlagSet
	fs.Init("test", zflag.ContinueOnError)
	fs.String("string1", "some", "string1 usage", zflag.OptShorthand('s'))
	fs.Bool("bool1", false, "bool1 usage", zflag.OptShorthand('b'))

	fs.String("string2", "some", "string2 usage in group1", zflag.OptGroup("group1"))
	fs.Bool("bool2", false, "bool2 usage in group1", zflag.OptGroup("group1"))

	fs.String("string3", "some", "string3 usage in group2", zflag.OptGroup("group2"))
	fs.Bool("bool3", false, "bool3 usage in group2", zflag.OptGroup("group2"))

	expectedGroupLen := 3
	if expectedGroupLen != len(fs.Groups()) {
		t.Fatalf("expected %d groups, got %d", expectedGroupLen, len(fs.Groups()))
	}

	expectedGroup1 := ""
	if fs.Groups()[0] != expectedGroup1 {
		t.Errorf("expected %q, got %q", expectedGroup1, fs.Groups()[0])
	}

	expectedGroup2 := "group1"
	if fs.Groups()[1] != expectedGroup2 {
		t.Errorf("expected %q, got %q", expectedGroup2, fs.Groups()[1])
	}

	expectedGroup3 := "group2"
	if fs.Groups()[2] != expectedGroup3 {
		t.Errorf("expected %q, got %q", expectedGroup3, fs.Groups()[2])
	}
}

func TestEmptyUngrouped(t *testing.T) {
	var fs zflag.FlagSet
	fs.Init("test", zflag.ContinueOnError)

	fs.String("string2", "some", "string2 usage in group1", zflag.OptGroup("group1"))
	fs.Bool("bool2", false, "bool2 usage in group1", zflag.OptGroup("group1"))

	fs.String("string3", "some", "string3 usage in group2", zflag.OptGroup("group2"))
	fs.Bool("bool3", false, "bool3 usage in group2", zflag.OptGroup("group2"))

	expectedGroupLen := 2
	if expectedGroupLen != len(fs.Groups()) {
		t.Fatalf("expected %d groups, got %d", expectedGroupLen, len(fs.Groups()))
	}

	expectedGroup1 := "group1"
	if fs.Groups()[0] != expectedGroup1 {
		t.Errorf("expected %q, got %q", expectedGroup1, fs.Groups()[0])
	}

	expectedGroup2 := "group2"
	if fs.Groups()[1] != expectedGroup2 {
		t.Errorf("expected %q, got %q", expectedGroup2, fs.Groups()[1])
	}
}

func TestUserDefined(t *testing.T) {
	var flags zflag.FlagSet
	flags.Init("test", zflag.ContinueOnError)
	var v flagVar
	flags.Var(&v, "v", "usage", zflag.OptShorthand('v'))
	if err := flags.Parse([]string{"--v=1", "-v2", "-v", "3"}); err != nil {
		t.Error(err)
	}
	if len(v) != 3 {
		t.Fatal("expected 3 args; got ", len(v))
	}
	expect := "[1 2 3]"
	if v.String() != expect {
		t.Errorf("expected value %q got %q", expect, v.String())
	}
}

func TestSetOutput(t *testing.T) {
	var flags zflag.FlagSet
	var buf strings.Builder
	flags.SetOutput(&buf)
	flags.Init("test", zflag.ContinueOnError)
	_ = flags.Parse([]string{"--unknown"})
	if out := buf.String(); !strings.Contains(out, "--unknown") {
		t.Fatalf("expected output mentioning unknown; got %q", out)
	}
}

func TestOutput(t *testing.T) {
	var flags zflag.FlagSet
	var buf strings.Builder
	expect := "an example string"
	flags.SetOutput(&buf)
	fmt.Fprint(flags.Output(), expect)
	if out := buf.String(); !strings.Contains(out, expect) {
		t.Fatalf("expected output %q; got %q", expect, out)
	}
}

func TestOutputExitOnError(t *testing.T) {
	if os.Getenv("ZFLAG_CRASH_TEST") == "1" {
		zflag.CommandLine = zflag.NewFlagSet(t.Name(), zflag.ExitOnError)
		os.Args = []string{t.Name(), "--unknown"}
		zflag.Parse()
		t.Fatal("this error should not be triggered")
		return
	}

	mockStdout := new(strings.Builder)
	mockStderr := new(strings.Builder)
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), "ZFLAG_CRASH_TEST=1")
	cmd.Stdout = mockStdout
	cmd.Stderr = mockStderr
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		want := "Usage of " + t.Name() + ":\n\nunknown flag: --unknown\n"
		if got := mockStderr.String(); got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
		if got := mockStdout.String(); len(got) != 0 {
			t.Errorf("stdout should be empty, got: %s", got)
		}
		return
	}
	t.Fatal("this error should not be triggered")
}

// This tests that one can reset the flags. This still works but not well, and is
// superseded by FlagSet.
func TestChangingArgs(t *testing.T) {
	zflag.ResetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--before", "subcmd"}
	before := zflag.Bool("before", false, "")
	if err := zflag.CommandLine.Parse(os.Args[1:]); err != nil {
		t.Fatal(err)
	}
	cmd := zflag.Arg(0)
	os.Args = []string{"subcmd", "--after", "args"}
	after := zflag.Bool("after", false, "")
	zflag.Parse()
	args := zflag.Args()

	if !*before || cmd != "subcmd" || !*after || len(args) != 1 || args[0] != "args" {
		t.Fatalf("expected true subcmd true [args] got %v %v %v %v", *before, cmd, *after, args)
	}
}

// Test that -help invokes the usage message and returns ErrHelp.
func TestHelp(t *testing.T) {
	var helpCalled = false
	fs := zflag.NewFlagSet("help test", zflag.ContinueOnError)
	fs.Usage = func() { helpCalled = true }

	var flag bool
	fs.BoolVar(&flag, "flag", false, "regular flag")

	// Regular flag invocation should work
	err := fs.Parse([]string{"--flag=true"})
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if !flag {
		t.Fatal("flag was not set by --flag")
	}
	if helpCalled {
		t.Fatal("help called for regular flag")
	}

	// Help flag should work as expected.
	for _, f := range []string{"--help", "-h", "-help", "-helpxyz", "-hxyz"} {
		err = fs.Parse([]string{f})
		if err == nil {
			t.Fatalf("while passing %s, error expected\n", f)
		}
		if err != zflag.ErrHelp {
			t.Fatalf("while passing %s, expected ErrHelp; got %s\n", f, err)
		}
		if !helpCalled {
			t.Fatalf("while passing %s, help was not called\n", f)
		}
		helpCalled = false
	}

	// Help flag should not work when disabled
	fs.DisableBuiltinHelp = true
	for _, f := range []string{"--help", "-h"} {
		err := fs.Parse([]string{f})
		if err == nil {
			t.Fatalf("while passing %s, error expected", f)
		}
		if err.Error() != "unknown flag: --help" && err.Error() != "unknown shorthand flag: 'h' in -h" {
			t.Fatalf("while passing %s, unknown flag error expected, got %s\n", f, err)
		}
		if !helpCalled {
			// Help should be triggered because this is an unknown, but not because the help flag was called
			t.Fatalf("while passing %s, help was not called\n", f)
		}
	}
	helpCalled = false
	// ... when disabled, any other shorthands should trigger an error
	for _, f := range []string{"-help", "-helpxyz", "-hxyz"} {
		err = fs.Parse([]string{f})
		if err == nil {
			t.Fatalf("while passing %s, error expected\n", f)
		}
	}
	helpCalled = false
	fs.ParseErrorsAllowList.UnknownFlags = true
	for _, f := range []string{"--help", "-h"} {
		err := fs.Parse([]string{f})
		t.Logf("help called: %v\n", helpCalled)
		if err != nil {
			t.Fatalf("while passing %s, error not expected, got %s\n", f, err)
		}
		if helpCalled {
			// Help should be triggered because this is an unknown, but not because the help flag was called
			t.Fatalf("while passing %s, help was not called\n", f)
		}
	}
	fs.DisableBuiltinHelp = false

	// If we define a help flag, that should override.
	var help bool
	fs.BoolVar(&help, "help", false, "help flag", zflag.OptShorthand('h'))
	err = fs.Parse([]string{"--help"})
	if err != nil {
		t.Fatal("expected no error for defined --help; got ", err)
	}
	if !help {
		t.Fatal("help should be true for defined --help")
	}
	if helpCalled {
		t.Fatal("help was called; should not have been for defined help flag")
	}
	help = false
	// ... including the shorthand
	err = fs.Parse([]string{"-h"})
	if err != nil {
		t.Fatal("expected no error for defined -h; got ", err)
	}
	if !help {
		t.Fatal("help should be true for defined -h")
	}
	if helpCalled {
		t.Fatal("help was called; should not have been for defined help flag")
	}
	help = false

	// If we define a help flag, that should override when the built in help flag is disabled.
	fs.DisableBuiltinHelp = true
	err = fs.Parse([]string{"--help"})
	if err != nil {
		t.Fatal("expected no error for defined --help; got ", err)
	}
	if !help {
		t.Fatal("help should be true for defined --help")
	}
	if helpCalled {
		t.Fatal("help was called; should not have been for defined help flag")
	}
	help = false
	// ... including the shorthand
	err = fs.Parse([]string{"-h"})
	if err != nil {
		t.Fatal("expected no error for defined -h; got ", err)
	}
	if !help {
		t.Fatal("help should be true for defined -h")
	}
	if helpCalled {
		t.Fatal("help was called; should not have been for defined help flag")
	}
}

func TestNoInterspersed(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetInterspersed(false)
	f.Bool("true", true, "always true")
	f.Bool("false", false, "always false")
	err := f.Parse([]string{"--true", "break", "--false"})
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	args := f.Args()
	if len(args) != 2 || args[0] != "break" || args[1] != "--false" {
		t.Fatal("expected interspersed options/non-options to fail")
	}
}

func TestTermination(t *testing.T) {
	f := zflag.NewFlagSet("termination", zflag.ContinueOnError)
	boolFlag := f.Bool("bool", false, "bool value", zflag.OptShorthand('l'))
	if f.Parsed() {
		t.Error("f.Parse() = true before Parse")
	}
	arg1 := "ls"
	arg2 := "-l"
	args := []string{
		"--",
		arg1,
		arg2,
	}
	f.SetOutput(ioutil.Discard)
	if err := f.Parse(args); err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if !f.Parsed() {
		t.Error("f.Parse() = false after Parse")
	}
	if *boolFlag {
		t.Error("expected boolFlag=false, got true")
	}
	if len(f.Args()) != 2 {
		t.Errorf("expected 2 arguments, got %d: %v", len(f.Args()), f.Args())
	}
	if f.Args()[0] != arg1 {
		t.Errorf("expected argument %q got %q", arg1, f.Args()[0])
	}
	if f.Args()[1] != arg2 {
		t.Errorf("expected argument %q got %q", arg2, f.Args()[1])
	}
	if f.ArgsLenAtDash() != 0 {
		t.Errorf("expected argsLenAtDash %d got %d", 0, f.ArgsLenAtDash())
	}
}

func getDeprecatedFlagSet() *zflag.FlagSet {
	f := zflag.NewFlagSet("bob", zflag.ContinueOnError)
	f.Bool("badflag", true, "always true", zflag.OptDeprecated("use --good-flag instead"))
	return f
}
func TestDeprecatedFlagInDocs(t *testing.T) {
	f := getDeprecatedFlagSet()

	out := new(strings.Builder)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "badflag") {
		t.Errorf("found deprecated flag in usage!")
	}
}

func TestUnHiddenDeprecatedFlagInDocs(t *testing.T) {
	f := getDeprecatedFlagSet()
	flg := f.Lookup("badflag")
	if flg == nil {
		t.Fatalf("Unable to lookup 'bob' in TestUnHiddenDeprecatedFlagInDocs")
	}
	flg.Hidden = false

	out := new(strings.Builder)
	f.SetOutput(out)
	f.PrintDefaults()

	defaults := out.String()
	if !strings.Contains(defaults, "badflag") {
		t.Errorf("Did not find deprecated flag in usage!")
	}
	if !strings.Contains(defaults, "use --good-flag instead") {
		t.Errorf("Did not find 'use --good-flag instead' in defaults")
	}
}

func TestDeprecatedFlagShorthandInDocs(t *testing.T) {
	f := zflag.NewFlagSet("bob", zflag.ContinueOnError)
	name := "noshorthandflag"
	f.Bool(name, true, "always true", zflag.OptShorthand('n'), zflag.OptShorthandDeprecated(fmt.Sprintf("use --%s instead", name)))

	out := new(strings.Builder)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "-n,") {
		t.Errorf("found deprecated flag shorthand in usage!")
	}
}

func parseReturnStderr(_ *testing.T, f *zflag.FlagSet, args []string) (string, error) {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := f.Parse(args)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf strings.Builder
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stderr = oldStderr
	out := <-outC

	return out, err
}

func TestDeprecatedFlagUsage(t *testing.T) {
	f := zflag.NewFlagSet("bob", zflag.ContinueOnError)
	usageMsg := "use --good-flag instead"
	f.Bool("badflag", true, "always true", zflag.OptDeprecated(usageMsg))

	args := []string{"--badflag"}
	out, err := parseReturnStderr(t, f, args)
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}

	if !strings.Contains(out, usageMsg) {
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	}
}

func TestDeprecatedFlagShorthandUsage(t *testing.T) {
	f := zflag.NewFlagSet("bob", zflag.ContinueOnError)
	name := "noshorthandflag"
	usageMsg := fmt.Sprintf("use --%s instead", name)
	f.Bool(name, true, "always true", zflag.OptShorthand('n'), zflag.OptShorthandDeprecated(usageMsg))

	args := []string{"-n"}
	out, err := parseReturnStderr(t, f, args)
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}

	if !strings.Contains(out, usageMsg) {
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	}
}

func TestDeprecatedFlagUsageNormalized(t *testing.T) {
	f := zflag.NewFlagSet("bob", zflag.ContinueOnError)
	usageMsg := "use --good-flag instead"
	f.Bool("bad-double_flag", true, "always true", zflag.OptDeprecated(usageMsg))
	f.SetNormalizeFunc(wordSepNormalizeFunc)

	args := []string{"--bad_double_flag"}
	out, err := parseReturnStderr(t, f, args)
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}

	if !strings.Contains(out, usageMsg) {
		t.Errorf("usageMsg not printed when using a deprecated flag!")
	}
}

// Name normalization function should be called only once on flag addition
func TestMultipleNormalizeFlagNameInvocations(t *testing.T) {
	normalizeFlagNameInvocations = 0

	f := zflag.NewFlagSet("normalized", zflag.ContinueOnError)
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	f.Bool("with_under_flag", false, "bool value")

	if normalizeFlagNameInvocations != 1 {
		t.Fatal("Expected normalizeFlagNameInvocations to be 1; got ", normalizeFlagNameInvocations)
	}
}

func TestHiddenFlagInUsage(t *testing.T) {
	f := zflag.NewFlagSet("bob", zflag.ContinueOnError)
	f.Bool("secretFlag", true, "shhh", zflag.OptHidden())

	out := new(strings.Builder)
	f.SetOutput(out)
	f.PrintDefaults()

	if strings.Contains(out.String(), "secretFlag") {
		t.Errorf("found hidden flag in usage!")
	}
}

func TestHiddenFlagUsage(t *testing.T) {
	f := zflag.NewFlagSet("bob", zflag.ContinueOnError)
	f.Bool("secretFlag", true, "shhh", zflag.OptHidden())

	args := []string{"--secretFlag"}
	out, err := parseReturnStderr(t, f, args)
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}

	if strings.Contains(out, "shhh") {
		t.Errorf("usage message printed when using a hidden flag!")
	}
}

const defaultOutput = `      --A                     for bootstrapping, allow 'any' type
      --[no-]Alongflagname    disable bounds checking
  -C, --[no-]CCC              a boolean defaulting to true (default true)
      --D path                set relative path for local imports
      --F number              a non-zero number (default 2.7)
      --G float               a float that defaults to zero
      --IP ip                 IP address with no default
      --IPMask ipMask         Netmask address with no default
      --IPNet ipNet           IP network with no default
      --Ints ints             int slice with zero default
      --N int                 a non-zero int (default 27)
      --StringSlice strings   string slice with zero default
      --Z int                 an int that defaults to zero
      --custom custom         custom Value implementation
      --customP custom        a VarP with default (default 10)
      --disableDefault int    A non-zero int with DisablePrintDefault
      --maxT timeout          set timeout for dial
  -v, --verbose count         verbosity
`

// Custom value that satisfies the Value interface.
type customValue int

func (cv *customValue) String() string { return fmt.Sprintf("%v", *cv) }

func (cv *customValue) Set(val string) error {
	v, err := strconv.ParseInt(val, 0, 64)
	*cv = customValue(v)
	return err
}

func (cv *customValue) Type() string { return "custom" }

func TestPrintDefaults(t *testing.T) {
	fs := zflag.NewFlagSet("print defaults test", zflag.ContinueOnError)
	var buf strings.Builder
	fs.SetOutput(&buf)
	fs.Bool("A", false, "for bootstrapping, allow 'any' type")
	fs.Bool("Alongflagname", false, "disable bounds checking", zflag.OptAddNegative())
	fs.Bool("CCC", true, "a boolean defaulting to true", zflag.OptShorthand('C'), zflag.OptAddNegative())
	fs.String("D", "", "set relative `path` for local imports")
	fs.Float64("F", 2.7, "a non-zero `number`")
	fs.Float64("G", 0, "a float that defaults to zero")
	fs.Int("N", 27, "a non-zero int")
	fs.IntSlice("Ints", []int{}, "int slice with zero default")
	fs.IP("IP", nil, "IP address with no default")
	fs.IPMask("IPMask", nil, "Netmask address with no default")
	fs.IPNet("IPNet", net.IPNet{}, "IP network with no default")
	fs.Int("Z", 0, "an int that defaults to zero")
	fs.Duration("maxT", 0, "set `timeout` for dial")
	fs.StringSlice("StringSlice", []string{}, "string slice with zero default")
	fs.Count("verbose", "verbosity", zflag.OptShorthand('v'))
	fs.Int("disableDefault", -1, "A non-zero int with DisablePrintDefault", zflag.OptDisablePrintDefault())

	var cv customValue
	fs.Var(&cv, "custom", "custom Value implementation")

	cv2 := customValue(10)
	fs.Var(&cv2, "customP", "a VarP with default")

	fs.PrintDefaults()
	got := buf.String()
	if got != defaultOutput {
		fmt.Printf("\n%s\n", got)
		fmt.Printf("\n%s\n", defaultOutput)
		t.Errorf("got %q want %q\n", got, defaultOutput)
	}
}

func TestVisitAllFlagOrder(t *testing.T) {
	fs := zflag.NewFlagSet("TestVisitAllFlagOrder", zflag.ContinueOnError)
	fs.SortFlags = false
	// https://github.com/spf13/zflag/issues/120
	fs.SetNormalizeFunc(func(_ *zflag.FlagSet, name string) zflag.NormalizedName {
		return zflag.NormalizedName(name)
	})

	names := []string{"C", "B", "A", "D"}
	for _, name := range names {
		fs.Bool(name, false, "")
	}

	i := 0
	fs.VisitAll(func(f *zflag.Flag) {
		if names[i] != f.Name {
			t.Errorf("Incorrect order. Expected %v, got %v", names[i], f.Name)
		}
		i++
	})
}

func TestVisitFlagOrder(t *testing.T) {
	fs := zflag.NewFlagSet("TestVisitFlagOrder", zflag.ContinueOnError)
	fs.SortFlags = false
	names := []string{"C", "B", "A", "D"}
	for _, name := range names {
		fs.Bool(name, false, "")
		_ = fs.Set(name, "true")
	}

	i := 0
	fs.Visit(func(f *zflag.Flag) {
		if names[i] != f.Name {
			t.Errorf("Incorrect order. Expected %v, got %v", names[i], f.Name)
		}
		i++
	})
}

func TestUnquoteUsage(t *testing.T) {
	tests := []struct {
		name                   string
		flagUsage              string
		expectedUsage          string
		opts                   []zflag.Opt
		overrideUnquoteUsage   bool
		disableUnquoteUsageVal bool
	}{
		{
			name:          "default unquotes",
			flagUsage:     "test `ctype1`",
			expectedUsage: "--test ctype1   test ctype1",
		},
		{
			name:                   "unquote when usage type set and unquote explicitly unset",
			flagUsage:              "test `ctype2`",
			opts:                   []zflag.Opt{zflag.OptUsageType("foo")},
			overrideUnquoteUsage:   true,
			disableUnquoteUsageVal: false,
			expectedUsage:          "--test foo   test ctype2",
		},
		{
			name:          "does not unquote when unquote usage disabled",
			flagUsage:     "test `ctype3`",
			opts:          []zflag.Opt{zflag.OptDisableUnquoteUsage()},
			expectedUsage: "--test string   test `ctype3`",
		},
		{
			name:          "disables unquote usage when usage type set",
			flagUsage:     "test `ctype4`",
			opts:          []zflag.Opt{zflag.OptUsageType("bar")},
			expectedUsage: "--test bar   test `ctype4`",
		},
		{
			name:          "Skips if single backtick",
			flagUsage:     "test `ctype4",
			expectedUsage: "--test string   test `ctype4",
		},
		{
			name:                   "Indexing start yes end no",
			flagUsage:              "`test ctype4",
			expectedUsage:          "--test string   `test ctype4",
			overrideUnquoteUsage:   true,
			disableUnquoteUsageVal: false,
		},
		{
			name:                   "Indexing start yes end no, custom usage type",
			flagUsage:              "`test ctype4",
			expectedUsage:          "--test val   `test ctype4",
			opts:                   []zflag.Opt{zflag.OptUsageType("val")},
			overrideUnquoteUsage:   true,
			disableUnquoteUsageVal: false,
		},
		{
			name:                   "Indexing start no end yes",
			flagUsage:              "`test ctype4",
			expectedUsage:          "--test string   `test ctype4",
			overrideUnquoteUsage:   true,
			disableUnquoteUsageVal: false,
		},
		{
			name:                   "Indexing start no end yes, custom usage type",
			flagUsage:              "`test ctype4",
			expectedUsage:          "--test val   `test ctype4",
			opts:                   []zflag.Opt{zflag.OptUsageType("val")},
			overrideUnquoteUsage:   true,
			disableUnquoteUsageVal: false,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf strings.Builder
			var fs = zflag.NewFlagSet("", zflag.ContinueOnError)
			fs.SetOutput(&buf)

			// normal usage
			fs.String("test", "", tt.flagUsage, tt.opts...)
			if tt.overrideUnquoteUsage {
				fs.Lookup("test").DisableUnquoteUsage = tt.disableUnquoteUsageVal
			}
			want := fmt.Sprintf("Usage:\n      %s\n", tt.expectedUsage)

			zflag.CallDefaultUsage(fs)
			assertEqual(t, want, buf.String())
		})
	}
}

// TestCustomFlagValue verifies that custom flag usage string doesn't change its "default" section after parsing
func TestCustomFlagDefValue(t *testing.T) {
	fs := zflag.NewFlagSet("TestCustomFlagDefValue", zflag.ContinueOnError)
	var buf strings.Builder
	fs.SetOutput(&buf)

	var cv customValue
	fs.Var(&cv, "customP", "a Var with no default")

	fs.PrintDefaults()
	beforeParse := buf.String()
	buf.Reset()

	args := []string{
		"--customP=10",
	}

	err := fs.Parse(args)
	assertNoErr(t, err)

	val := fs.Lookup("customP").Value.String()
	assertEqual(t, "10", val)

	fs.PrintDefaults()
	afterParse := buf.String()

	if beforeParse != afterParse {
		fmt.Print("\n" + beforeParse + "\n")
		fmt.Print("\n" + afterParse + "\n")
		t.Errorf("got %q want %q\n", afterParse, beforeParse)
	}
}

// TestNoDuplicateUnknownFlagError ensures issue https://github.com/spf13/pflag/issues/352 does not regress.
func TestNoDuplicateUnknownFlagError(t *testing.T) {
	zflag.SetExitFunc(func(code int) {
		assertEqual(t, 2, code)
	})

	var buf strings.Builder

	fs := zflag.NewFlagSet("myprog", zflag.ExitOnError)
	fs.SetOutput(&buf)
	err := fs.Parse([]string{"--bogus"})
	assertNoErr(t, err)

	substr := "unknown flag: --bogus"
	count := strings.Count(buf.String(), substr)
	assertEqualf(t, 1, count, "expected %q to appear in output exactly once, got %d", substr, count)
}
