// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

// -- string Value
type stringValue string

var _ Value = (*stringValue)(nil)
var _ Getter = (*stringValue)(nil)
var _ Typed = (*stringValue)(nil)

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() any {
	return string(*s)
}

func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string { return string(*s) }

// GetString return the string value of a flag with the given name
func (fs *FlagSet) GetString(name string) (string, error) {
	val, err := fs.getFlagValue(name, "string")
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

// MustGetString is like GetString, but panics on error.
func (fs *FlagSet) MustGetString(name string) string {
	val, err := fs.GetString(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (fs *FlagSet) StringVar(p *string, name string, value string, usage string, opts ...Opt) {
	fs.Var(newStringValue(value, p), name, usage, opts...)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringVar(p *string, name string, value string, usage string, opts ...Opt) {
	CommandLine.StringVar(p, name, value, usage, opts...)
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (fs *FlagSet) String(name string, value string, usage string, opts ...Opt) *string {
	var p string
	fs.StringVar(&p, name, value, usage, opts...)
	return &p
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, value string, usage string, opts ...Opt) *string {
	return CommandLine.String(name, value, usage, opts...)
}
