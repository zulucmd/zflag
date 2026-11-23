// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// -- net.IPMask value
type ipMaskValue net.IPMask

var _ Value = (*ipMaskValue)(nil)
var _ Getter = (*ipMaskValue)(nil)
var _ Typed = (*ipMaskValue)(nil)

func newIPMaskValue(val net.IPMask, p *net.IPMask) *ipMaskValue {
	*p = val
	return (*ipMaskValue)(p)
}

func (i *ipMaskValue) String() string { return net.IPMask(*i).String() }
func (i *ipMaskValue) Set(val string) error {
	val = strings.TrimSpace(val)
	ip := ParseIPv4Mask(val)
	if ip == nil {
		return fmt.Errorf("failed to parse IP mask: %q", val)
	}
	*i = ipMaskValue(ip)
	return nil
}

func (i *ipMaskValue) Get() any {
	return net.IPMask(*i)
}

func (i *ipMaskValue) Type() string {
	return "ipMask"
}

// ParseIPv4Mask written in IP form (e.g. 255.255.255.0).
// This function should really belong to the net package.
func ParseIPv4Mask(s string) net.IPMask {
	mask := net.ParseIP(s)
	if mask == nil {
		if len(s) != 8 {
			return nil
		}
		// net.IPMask.String() actually outputs things like ffffff00
		// so write a horrible parser for that as well  :-(
		m := []int{}
		for i := 0; i < 4; i++ {
			b := "0x" + s[2*i:2*i+2]
			d, err := strconv.ParseInt(b, 0, 0)
			if err != nil {
				return nil
			}
			m = append(m, int(d))
		}
		s := fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
		mask = net.ParseIP(s)
		if mask == nil {
			return nil
		}
	}
	return net.IPv4Mask(mask[12], mask[13], mask[14], mask[15])
}

// GetIPMask return the net.IPv4Mask value of a flag with the given name
func (fs *FlagSet) GetIPMask(name string) (net.IPMask, error) {
	val, err := fs.getFlagValue(name, "ipMask")
	if err != nil {
		return nil, err
	}
	return val.(net.IPMask), nil
}

// MustGetIPMask is like GetIPMask, but panics on error.
func (fs *FlagSet) MustGetIPMask(name string) net.IPMask {
	val, err := fs.GetIPMask(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IPMaskVar defines a net.IPMask flag with specified name, default value, and usage string.
// The argument p points to a net.IPMask variable in which to store the value of the flag.
func (fs *FlagSet) IPMaskVar(p *net.IPMask, name string, value net.IPMask, usage string, opts ...Opt) {
	fs.Var(newIPMaskValue(value, p), name, usage, opts...)
}

// IPMaskVar defines a net.IPMask flag with specified name, default value, and usage string.
// The argument p points to a net.IPMask variable in which to store the value of the flag.
func IPMaskVar(p *net.IPMask, name string, value net.IPMask, usage string, opts ...Opt) {
	CommandLine.IPMaskVar(p, name, value, usage, opts...)
}

// IPMask defines a net.IPMask flag with specified name, default value, and usage string.
// The return value is the address of a net.IPMask variable that stores the value of the flag.
func (fs *FlagSet) IPMask(name string, value net.IPMask, usage string, opts ...Opt) *net.IPMask {
	var p net.IPMask
	fs.IPMaskVar(&p, name, value, usage, opts...)
	return &p
}

// IPMask defines a net.IPMask flag with specified name, default value, and usage string.
// The return value is the address of a net.IPMask variable that stores the value of the flag.
func IPMask(name string, value net.IPMask, usage string, opts ...Opt) *net.IPMask {
	return CommandLine.IPMask(name, value, usage, opts...)
}
