// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// BytesHex adapts []byte for use as a flag. Value of flag is HEX encoded
type bytesHexValue []byte

var _ Value = (*bytesHexValue)(nil)
var _ Getter = (*bytesHexValue)(nil)
var _ Typed = (*bytesHexValue)(nil)

// String implements zflag.Value.
func (bytesHex *bytesHexValue) String() string {
	return fmt.Sprintf("%X", *bytesHex)
}

func (bytesHex *bytesHexValue) Get() any {
	return []byte(*bytesHex)
}

// Set implements zflag.Value.Set.
func (bytesHex *bytesHexValue) Set(value string) error {
	bin, err := hex.DecodeString(strings.TrimSpace(value))

	if err != nil {
		return err
	}

	*bytesHex = bin

	return nil
}

// Type implements zflag.Value.Type.
func (*bytesHexValue) Type() string {
	return "bytesHex"
}

func newBytesHexValue(val []byte, p *[]byte) *bytesHexValue {
	*p = val
	return (*bytesHexValue)(p)
}

// GetBytesHex return the []byte value of a flag with the given name
func (fs *FlagSet) GetBytesHex(name string) ([]byte, error) {
	val, err := fs.getFlagValue(name, "bytesHex")

	if err != nil {
		return []byte{}, err
	}

	return val.([]byte), nil
}

// MustGetBytesHex is like GetBytesHex, but panics on error.
func (fs *FlagSet) MustGetBytesHex(name string) []byte {
	val, err := fs.GetBytesHex(name)
	if err != nil {
		panic(err)
	}
	return val
}

// BytesHexVar defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func (fs *FlagSet) BytesHexVar(p *[]byte, name string, value []byte, usage string, opts ...Opt) {
	fs.Var(newBytesHexValue(value, p), name, usage, opts...)
}

// BytesHexVar defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func BytesHexVar(p *[]byte, name string, value []byte, usage string, opts ...Opt) {
	CommandLine.BytesHexVar(p, name, value, usage, opts...)
}

// BytesHex defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func (fs *FlagSet) BytesHex(name string, value []byte, usage string, opts ...Opt) *[]byte {
	var p []byte
	fs.BytesHexVar(&p, name, value, usage, opts...)
	return &p
}

// BytesHex defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func BytesHex(name string, value []byte, usage string, opts ...Opt) *[]byte {
	return CommandLine.BytesHex(name, value, usage, opts...)
}

// BytesBase64 adapts []byte for use as a flag. Value of flag is Base64 encoded
type bytesBase64Value []byte

var _ Value = (*bytesBase64Value)(nil)
var _ Getter = (*bytesBase64Value)(nil)
var _ Typed = (*bytesBase64Value)(nil)

// String implements zflag.Value.String.
func (bytesBase64 *bytesBase64Value) String() string {
	return base64.StdEncoding.EncodeToString(*bytesBase64)
}

func (bytesBase64 *bytesBase64Value) Get() any {
	return []byte(*bytesBase64)
}

// Set implements zflag.Value.Set.
func (bytesBase64 *bytesBase64Value) Set(value string) error {
	bin, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return err
	}

	*bytesBase64 = bin

	return nil
}

// Type implements zflag.Value.Type.
func (*bytesBase64Value) Type() string {
	return "bytesBase64"
}

func newBytesBase64Value(val []byte, p *[]byte) *bytesBase64Value {
	*p = val
	return (*bytesBase64Value)(p)
}

// GetBytesBase64 return the []byte value of a flag with the given name
func (fs *FlagSet) GetBytesBase64(name string) ([]byte, error) {
	val, err := fs.getFlagValue(name, "bytesBase64")
	if err != nil {
		return []byte{}, err
	}
	return val.([]byte), nil
}

// MustGetBytesBase64 is like GetBytesBase64, but panics on error.
func (fs *FlagSet) MustGetBytesBase64(name string) []byte {
	val, err := fs.GetBytesBase64(name)
	if err != nil {
		panic(err)
	}
	return val
}

// BytesBase64Var defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func (fs *FlagSet) BytesBase64Var(p *[]byte, name string, value []byte, usage string, opts ...Opt) {
	fs.Var(newBytesBase64Value(value, p), name, usage, opts...)
}

// BytesBase64Var defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func BytesBase64Var(p *[]byte, name string, value []byte, usage string, opts ...Opt) {
	CommandLine.BytesBase64Var(p, name, value, usage, opts...)
}

// BytesBase64 defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func (fs *FlagSet) BytesBase64(name string, value []byte, usage string, opts ...Opt) *[]byte {
	var p []byte
	fs.BytesBase64Var(&p, name, value, usage, opts...)
	return &p
}

// BytesBase64 defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func BytesBase64(name string, value []byte, usage string, opts ...Opt) *[]byte {
	return CommandLine.BytesBase64(name, value, usage, opts...)
}
