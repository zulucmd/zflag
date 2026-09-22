// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"encoding"
	"fmt"
	"reflect"
)

// -- text Value
type textValue struct{ p encoding.TextUnmarshaler }

var _ Value = (*textValue)(nil)
var _ Getter = (*textValue)(nil)
var _ Typed = (*textValue)(nil)

func newTextValue(val encoding.TextMarshaler, p encoding.TextUnmarshaler) *textValue {
	ptrVal := reflect.ValueOf(p)
	if ptrVal.Kind() != reflect.Ptr {
		panic("variable value type must be a pointer")
	}
	defVal := reflect.ValueOf(val)
	if defVal.Kind() == reflect.Ptr {
		defVal = defVal.Elem()
	}
	if defVal.Type() != ptrVal.Type().Elem() {
		panic(fmt.Sprintf("default type does not match variable type: %v != %v", defVal.Type(), ptrVal.Type().Elem()))
	}
	ptrVal.Elem().Set(defVal)
	return &textValue{p}
}

func (t *textValue) Set(val string) error {
	return t.p.UnmarshalText([]byte(val))
}

func (t *textValue) Get() any {
	return t.p
}

func (t *textValue) String() string {
	if m, ok := t.p.(encoding.TextMarshaler); ok {
		if b, err := m.MarshalText(); err == nil {
			return string(b)
		}
	}
	return ""
}

func (t *textValue) Type() string {
	return "text"
}

// GetText sets out, which must implement encoding.TextUnmarshaler, to the
// value of the flag with the given name.
func (fs *FlagSet) GetText(name string, out encoding.TextUnmarshaler) error {
	flag := fs.Lookup(name)
	if flag == nil {
		return NewUnknownFlagError(name)
	}
	if v, isTyped := flag.Value.(Typed); isTyped && v.Type() != "text" {
		return fmt.Errorf("trying to get %q value of flag of type %q", "text", v.Type())
	}
	return out.UnmarshalText([]byte(flag.Value.String()))
}

// MustGetText is like GetText, but panics on error.
func (fs *FlagSet) MustGetText(name string, out encoding.TextUnmarshaler) {
	if err := fs.GetText(name, out); err != nil {
		panic(err)
	}
}

// Text defines a text flag with specified name, default value, and usage
// string. The return value is the unmarshal target that stores the value of
// the flag.
func (fs *FlagSet) Text(name string, value encoding.TextMarshaler, usage string, opts ...Opt) encoding.TextUnmarshaler {
	valueType := reflect.TypeOf(value)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	out, ok := reflect.New(valueType).Interface().(encoding.TextUnmarshaler)
	if !ok {
		panic(fmt.Sprintf("type %T does not implement encoding.TextUnmarshaler", value))
	}
	fs.TextVar(out, name, value, usage, opts...)
	return out
}

// Text defines a text flag with specified name, default value, and usage
// string. The return value is the unmarshal target that stores the value of
// the flag.
func Text(name string, value encoding.TextMarshaler, usage string, opts ...Opt) encoding.TextUnmarshaler {
	return CommandLine.Text(name, value, usage, opts...)
}

// TextVar defines a flag with given name, default value, and usage. p must be
// a pointer to a type implementing encoding.TextUnmarshaler; value must be of
// the same type as p.
func (fs *FlagSet) TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string, opts ...Opt) {
	fs.Var(newTextValue(value, p), name, usage, opts...)
}

// TextVar defines a flag with given name, default value, and usage. p must be
// a pointer to a type implementing encoding.TextUnmarshaler; value must be of
// the same type as p.
func TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string, opts ...Opt) {
	CommandLine.Var(newTextValue(value, p), name, usage, opts...)
}
