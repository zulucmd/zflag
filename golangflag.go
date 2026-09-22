// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	goflag "flag"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

// go test flags prefixes
func isGotestFlag(flag string) bool {
	return strings.HasPrefix(flag, "-test.")
}

func isGotestShorthandFlag(flag string) bool {
	return strings.HasPrefix(flag, "test.")
}

// flagValueWrapper implements zflag.Value around a flag.Value.  The main
// difference here is the addition of the Type method that returns a string
// name of the type.  As this is generally unknown, we approximate that with
// reflection.
type flagValueWrapper struct {
	inner    goflag.Value
	flagType string
}

var _ Value = (*flagValueWrapper)(nil)
var _ Getter = (*flagValueWrapper)(nil)
var _ Typed = (*flagValueWrapper)(nil)

func wrapFlagValue(v goflag.Value) Value {
	// If the flag.Value happens to also be a zflag.Value, just use it directly.
	if pv, ok := v.(Value); ok {
		return pv
	}

	pv := &flagValueWrapper{
		inner: v,
	}

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Interface || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	pv.flagType = strings.TrimSuffix(t.Name(), "Value")
	return pv
}

func (v *flagValueWrapper) String() string {
	return v.inner.String()
}

func (v *flagValueWrapper) Set(val string) error {
	return v.inner.Set(val)
}

func (v *flagValueWrapper) Get() interface{} {
	if getter, ok := v.inner.(goflag.Getter); ok {
		return getter.Get()
	}

	return v.inner.String()
}

func (v *flagValueWrapper) Type() string {
	return v.flagType
}

// FromGoFlag will return a *zflag.Flag given a *flag.Flag
// If the *flag.Flag.Name was a single character (ex: `v`) it will be accessible
// with both `-v` and `--v` in flags. If the golang flag was more than a single
// character (ex: `verbose`) it will only be accessible via `--verbose`
func FromGoFlag(goflag *goflag.Flag) *Flag {
	// Remember the default value as a string; it won't change.
	flag := &Flag{
		Name:  goflag.Name,
		Usage: goflag.Usage,
		Value: wrapFlagValue(goflag.Value),
		// Looks like golang flags don't set DefValue correctly  :-(
		// DefValue: goflag.DefValue,
		DefValue: goflag.Value.String(),
	}
	// Ex: if the golang flag was -v, allow both -v and --v to work
	if utf8.RuneCountInString(flag.Name) == 1 {
		short, _ := utf8.DecodeRuneInString(flag.Name)
		flag.Shorthand = short
	}
	return flag
}

// AddGoFlag will add the given *flag.Flag to the zflag.FlagSet
func (fs *FlagSet) AddGoFlag(goflag *goflag.Flag) {
	if fs.Lookup(goflag.Name) != nil {
		return
	}
	newflag := FromGoFlag(goflag)
	fs.AddFlag(newflag)
}

// AddGoFlagSet will add the given *flag.FlagSet to the zflag.FlagSet
func (fs *FlagSet) AddGoFlagSet(newSet *goflag.FlagSet) {
	if newSet == nil {
		return
	}
	newSet.VisitAll(func(goflag *goflag.Flag) {
		fs.AddGoFlag(goflag)
	})
	if fs.addedGoFlagSets == nil {
		fs.addedGoFlagSets = make([]*goflag.FlagSet, 0)
	}
	fs.addedGoFlagSets = append(fs.addedGoFlagSets, newSet)
}

// CopyToGoFlagSet adds all current flags to the given Go flag set.
// Deprecation remarks get copied into the usage description.
// Whenever possible, a flag gets added for which Go flags shows
// a proper type in the help message.
func (fs *FlagSet) CopyToGoFlagSet(newSet *goflag.FlagSet) {
	fs.VisitAll(func(flag *Flag) {
		usage := flag.Usage
		if flag.Deprecated != "" {
			usage += " (DEPRECATED: " + flag.Deprecated + ")"
		}

		switch value := flag.Value.(type) {
		case *stringValue:
			newSet.StringVar((*string)(value), flag.Name, *(*string)(value), usage)
		case *intValue:
			newSet.IntVar((*int)(value), flag.Name, *(*int)(value), usage)
		case *int64Value:
			newSet.Int64Var((*int64)(value), flag.Name, *(*int64)(value), usage)
		case *uintValue:
			newSet.UintVar((*uint)(value), flag.Name, *(*uint)(value), usage)
		case *uint64Value:
			newSet.Uint64Var((*uint64)(value), flag.Name, *(*uint64)(value), usage)
		case *durationValue:
			newSet.DurationVar((*time.Duration)(value), flag.Name, *(*time.Duration)(value), usage)
		case *float64Value:
			newSet.Float64Var((*float64)(value), flag.Name, *(*float64)(value), usage)
		default:
			newSet.Var(flag.Value, flag.Name, usage)
		}
	})
}

// ParseSkippedFlags parses go test flags (the ones starting with "-test.") with the
// given goflag.FlagSet, since zflag.Parse() skips them.
// Typical usage: ParseSkippedFlags(os.Args[1:], goflag.CommandLine)
func ParseSkippedFlags(osArgs []string, goFlagSet *goflag.FlagSet) error {
	var skippedFlags []string
	for _, f := range osArgs {
		if isGotestFlag(f) {
			skippedFlags = append(skippedFlags, f)
		}
	}
	return goFlagSet.Parse(skippedFlags)
}
