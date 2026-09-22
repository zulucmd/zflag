// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

// -- boolfunc Value
type boolFuncValue func(string) error

var _ Value = (*boolFuncValue)(nil)
var _ BoolFlag = (*boolFuncValue)(nil)
var _ Typed = (*boolFuncValue)(nil)

func newBoolFuncValue(fn func(string) error) *boolFuncValue {
	boolFuncVal := boolFuncValue(fn)
	return &boolFuncVal
}

func (i *boolFuncValue) Set(val string) error {
	return (*i)(val)
}

func (i *boolFuncValue) Type() string {
	return "func"
}

func (i *boolFuncValue) String() string { return "" }

func (i *boolFuncValue) IsBoolFlag() bool { return true }

// BoolFunc defines a flag with specified name, and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func (fs *FlagSet) BoolFunc(name string, usage string, fn func(string) error, opts ...Opt) {
	fs.Var(newBoolFuncValue(fn), name, usage, opts...)
}

// BoolFunc defines a flag with specified name, and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func BoolFunc(name string, usage string, fn func(string) error, opts ...Opt) {
	CommandLine.BoolFunc(name, usage, fn, opts...)
}

// These are not needed for this specific type, and they are added here to stop validate_types.sh from fail.
// func (fs *FlagSet) GetBoolFunc(
// func (fs *FlagSet) MustGetBoolFunc(
// func (fs *FlagSet) BoolFuncVar(
// func BoolFuncVar(
