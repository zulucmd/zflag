// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

// -- int Value
type funcValue func(string) error

var _ Value = (*funcValue)(nil)
var _ Typed = (*funcValue)(nil)

func newFuncValue(fn func(string) error) *funcValue {
	funcVal := funcValue(fn)
	return &funcVal
}

func (i *funcValue) Set(val string) error {
	return (*i)(val)
}

func (i *funcValue) Type() string {
	return "func"
}

func (i *funcValue) String() string { return "" }

// Func defines a flag with specified name, and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func (fs *FlagSet) Func(name string, usage string, fn func(string) error, opts ...Opt) {
	fs.Var(newFuncValue(fn), name, usage, opts...)
}

// Func defines a flag with specified name, and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func Func(name string, usage string, fn func(string) error, opts ...Opt) {
	CommandLine.Func(name, usage, fn, opts...)
}

// These are not needed for this specific type, and they are added here to stop validate_funcs.sh from fail.
// func (f *FlagSet) GetFunc(
// func (f *FlagSet) MustGetFunc(
// func (f *FlagSet) FuncVar(
// func FuncVar(
