// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// -- stringToString Value
type stringToStringValue struct {
	value         *map[string]string
	changed       bool
	valueOptional bool
}

var _ Value = (*stringToStringValue)(nil)
var _ Getter = (*stringToStringValue)(nil)
var _ Typed = (*stringToStringValue)(nil)
var _ MapValue = (*stringToStringValue)(nil)

func newStringToStringValue(val map[string]string, p *map[string]string) *stringToStringValue {
	ssv := new(stringToStringValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringToStringValue) Set(val string) error {
	kv := strings.SplitN(val, "=", 2)
	if !s.valueOptional && len(kv) != 2 {
		return fmt.Errorf("%q must be formatted as key=value", val)
	}

	key := kv[0]
	val = ""
	if len(kv) == 2 {
		val = kv[1]
	}

	if !s.changed {
		*s.value = map[string]string{}
	}

	(*s.value)[key] = val
	s.changed = true

	return nil
}

func (s *stringToStringValue) Get() interface{} {
	return *s.value
}

func (s *stringToStringValue) Type() string {
	return "stringToString"
}

func (s *stringToStringValue) IsMap() bool { return true }

func (s *stringToStringValue) String() string {
	keys := make([]string, 0, len(*s.value))
	for k := range *s.value {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	records := make([]string, 0, len(*s.value)>>1)
	for _, k := range keys {
		records = append(records, k+"="+strconv.Quote((*s.value)[k]))
	}

	return fmt.Sprintf("%s", records)
}

// GetStringToString return the map[string]string value of a flag with the given name
func (fs *FlagSet) GetStringToString(name string) (map[string]string, error) {
	val, err := fs.getFlagValue(name, "stringToString")
	if err != nil {
		return map[string]string{}, err
	}
	return val.(map[string]string), nil
}

// MustGetStringToString is like GetStringToString, but panics on error.
func (fs *FlagSet) MustGetStringToString(name string) map[string]string {
	val, err := fs.GetStringToString(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringToStringVar defines a map[string]string flag with specified name, default value, and usage string.
// The argument p points to a map[string]string variable in which to store the values of multiple flags.
func (fs *FlagSet) StringToStringVar(p *map[string]string, name string, value map[string]string, usage string, opts ...Opt) {
	fs.Var(newStringToStringValue(value, p), name, usage, opts...)
}

// StringToStringVar defines a map[string]string flag with specified name, default value, and usage string.
// The argument p points to a map[string]string variable in which to store the values of multiple flags.
func StringToStringVar(p *map[string]string, name string, value map[string]string, usage string, opts ...Opt) {
	CommandLine.StringToStringVar(p, name, value, usage, opts...)
}

// StringToString defines a map[string]string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]string variable that stores the values of multiple flags.
func (fs *FlagSet) StringToString(name string, value map[string]string, usage string, opts ...Opt) *map[string]string {
	var p map[string]string
	fs.StringToStringVar(&p, name, value, usage, opts...)
	return &p
}

// StringToString defines a map[string]string flag with specified name, default value, and usage string.
// The return value is the address of a map[string]string variable that stores the values of multiple flags.
func StringToString(name string, value map[string]string, usage string, opts ...Opt) *map[string]string {
	return CommandLine.StringToString(name, value, usage, opts...)
}

// OptMapValueOptional allows map flags to be set with a key alone, without the
// key=value form. It applies to string-to-string, string-to-int, and
// string-to-int64 flags.
func OptMapValueOptional() Opt {
	return func(f *Flag) error {
		switch v := f.Value.(type) {
		case *stringToStringValue:
			v.valueOptional = true
			return nil
		case *stringToIntValue:
			v.valueOptional = true
			return nil
		case *stringToInt64Value:
			v.valueOptional = true
			return nil
		}

		return fmt.Errorf("value of type %T cannot be optional", f.Value)
	}
}
