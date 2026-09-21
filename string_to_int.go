// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- stringToInt Value
type stringToIntValue struct {
	value         *map[string]int
	changed       bool
	valueOptional bool
}

var _ Value = (*stringToIntValue)(nil)
var _ Getter = (*stringToIntValue)(nil)
var _ Typed = (*stringToIntValue)(nil)
var _ MapValue = (*stringToIntValue)(nil)

func newStringToIntValue(val map[string]int, p *map[string]int) *stringToIntValue {
	ssv := new(stringToIntValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

// Format: a=1
func (s *stringToIntValue) Set(val string) error {
	kv := strings.SplitN(val, "=", 2)
	if len(kv) != 2 {
		return fmt.Errorf("%s must be formatted as key=value", val)
	}
	key, val := kv[0], kv[1]

	val = strings.TrimSpace(val)
	v, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	if !s.changed {
		*s.value = map[string]int{}
	}

	(*s.value)[key] = v
	s.changed = true

	return nil
}

func (s *stringToIntValue) Get() interface{} {
	return *s.value
}

func (s *stringToIntValue) Type() string {
	return "stringToInt"
}

func (s *stringToIntValue) IsMap() bool { return true }

func (s *stringToIntValue) String() string {
	records := make([]string, 0, len(*s.value)>>1)
	for k, v := range *s.value {
		records = append(records, k+"="+strconv.Itoa(v))
	}

	return fmt.Sprintf("%s", records)
}

// GetStringToInt return the map[string]int value of a flag with the given name
func (fs *FlagSet) GetStringToInt(name string) (map[string]int, error) {
	val, err := fs.getFlagValue(name, "stringToInt")
	if err != nil {
		return map[string]int{}, err
	}
	return val.(map[string]int), nil
}

// MustGetStringToInt is like GetStringToInt, but panics on error.
func (fs *FlagSet) MustGetStringToInt(name string) map[string]int {
	val, err := fs.GetStringToInt(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringToIntVar defines a map[string]int flag with specified name, default value, and usage string.
// The argument p points to a map[string]int variable in which to store the values of multiple flags.
func (fs *FlagSet) StringToIntVar(p *map[string]int, name string, value map[string]int, usage string, opts ...Opt) {
	fs.Var(newStringToIntValue(value, p), name, usage, opts...)
}

// StringToIntVar defines a map[string]int flag with specified name, default value, and usage string.
// The argument p points to a map[string]int variable in which to store the values of multiple flags.
func StringToIntVar(p *map[string]int, name string, value map[string]int, usage string, opts ...Opt) {
	CommandLine.StringToIntVar(p, name, value, usage, opts...)
}

// StringToInt defines a map[string]int flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int variable that stores the values of multiple flags.
func (fs *FlagSet) StringToInt(name string, value map[string]int, usage string, opts ...Opt) *map[string]int {
	var p map[string]int
	fs.StringToIntVar(&p, name, value, usage, opts...)
	return &p
}

// StringToInt defines a map[string]int flag with specified name, default value, and usage string.
// The return value is the address of a map[string]int variable that stores the values of multiple flags.
func StringToInt(name string, value map[string]int, usage string, opts ...Opt) *map[string]int {
	return CommandLine.StringToInt(name, value, usage, opts...)
}
