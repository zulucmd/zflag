// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"net"
	"strings"
)

// IPNet adapts net.IPNet for use as a flag.
type ipNetValue net.IPNet

var _ Value = (*ipNetValue)(nil)
var _ Getter = (*ipNetValue)(nil)
var _ Typed = (*ipNetValue)(nil)

func (ipnet ipNetValue) String() string {
	n := net.IPNet(ipnet)
	return n.String()
}

func (ipnet *ipNetValue) Get() any {
	return net.IPNet(*ipnet)
}

func (ipnet *ipNetValue) Set(value string) error {
	_, n, err := net.ParseCIDR(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	*ipnet = ipNetValue(*n)
	return nil
}

func (*ipNetValue) Type() string {
	return "ipNet"
}

func newIPNetValue(val net.IPNet, p *net.IPNet) *ipNetValue {
	*p = val
	return (*ipNetValue)(p)
}

// GetIPNet return the net.IPNet value of a flag with the given name
func (fs *FlagSet) GetIPNet(name string) (net.IPNet, error) {
	val, err := fs.getFlagValue(name, "ipNet")
	if err != nil {
		return net.IPNet{}, err
	}
	return val.(net.IPNet), nil
}

// MustGetIPNet is like GetIPNet, but panics on error.
func (fs *FlagSet) MustGetIPNet(name string) net.IPNet {
	val, err := fs.GetIPNet(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IPNetVar defines a net.IPNet flag with specified name, default value, and usage string.
// The argument p points to a net.IPNet variable in which to store the value of the flag.
func (fs *FlagSet) IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string, opts ...Opt) {
	fs.Var(newIPNetValue(value, p), name, usage, opts...)
}

// IPNetVar defines a net.IPNet flag with specified name, default value, and usage string.
// The argument p points to a net.IPNet variable in which to store the value of the flag.
func IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string, opts ...Opt) {
	CommandLine.IPNetVar(p, name, value, usage, opts...)
}

// IPNet defines a net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of a net.IPNet variable that stores the value of the flag.
func (fs *FlagSet) IPNet(name string, value net.IPNet, usage string, opts ...Opt) *net.IPNet {
	var p net.IPNet
	fs.IPNetVar(&p, name, value, usage, opts...)
	return &p
}

// IPNet defines a net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of a net.IPNet variable that stores the value of the flag.
func IPNet(name string, value net.IPNet, usage string, opts ...Opt) *net.IPNet {
	return CommandLine.IPNet(name, value, usage, opts...)
}
