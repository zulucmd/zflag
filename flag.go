// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ErrHelp is the error returned if the flag -help is invoked but no such flag is defined.
var ErrHelp = errors.New("zflag: help requested")

// ErrorHandling defines how to handle flag parsing errors.
type ErrorHandling int

const (
	// ContinueOnError will return an err from Parse() if an error is found
	ContinueOnError ErrorHandling = iota
	// ExitOnError will call os.Exit(2) if an error is found when parsing
	ExitOnError
	// PanicOnError will panic() if an error is found when parsing flags
	PanicOnError
)

// ParseErrorsAllowList defines the parsing errors that can be ignored
type ParseErrorsAllowList struct {
	// UnknownFlags will ignore unknown flags errors and continue parsing rest of the flags
	// See GetUnknownFlags to retrieve collected unknowns.
	UnknownFlags bool
	// RequiredFlags will ignore required flags errors and continue parsing rest of the flags
	// See GetRequiredFlags to retrieve collected required flags.
	RequiredFlags bool
}

// NormalizedName is a flag name that has been normalized according to rules
// for the FlagSet (e.g. making '-' and '_' equivalent).
type NormalizedName string

// A FlagSet represents a set of defined flags.
type FlagSet struct {
	// Usage is the function called when an error occurs while parsing flags.
	// The field is a function (not a method) that may be changed to point to
	// a custom error handler.
	Usage func()

	// SortFlags is used to indicate, if user wants to have sorted flags in
	// help/usage messages.
	SortFlags bool

	// ParseErrorsAllowList is used to configure an allowlist of errors
	ParseErrorsAllowList ParseErrorsAllowList

	// DisableBuiltinHelp toggles the built-in convention of handling -h and --help
	DisableBuiltinHelp bool

	// FlagUsageFormatter allows for custom formatting of flag usage output.
	// Each individual item needs to be implemented. See FlagUsagesForGroupWrapped for info on what gets passed.
	FlagUsageFormatter FlagUsageFormatter

	name              string
	parsed            bool
	actual            map[NormalizedName]*Flag
	orderedActual     []*Flag
	sortedActual      []*Flag
	formal            map[NormalizedName]*Flag
	orderedFormal     []*Flag
	sortedFormal      []*Flag
	shorthands        map[rune]*Flag
	args              []string // arguments after flags
	argsLenAtDash     int      // len(args) when a '--' was located when parsing, or -1 if no --
	errorHandling     ErrorHandling
	output            io.Writer // nil means stderr; use Output() accessor
	interspersed      bool      // Allow interspersed option/non-option args
	normalizeNameFunc func(f *FlagSet, name string) NormalizedName

	addedGoFlagSets []*goflag.FlagSet
	unknownFlags    []string
}

// A Flag represents the state of a flag.
type Flag struct {
	Name                string              // Name as it appears on command line.
	Shorthand           rune                // Shorthand represents a one-letter abbreviation of a flag.
	ShorthandOnly       bool                // ShorthandOnly specifies if the user set only the shorthand.
	Usage               string              // Usage should contain the help message.
	UsageType           string              // UsageType is the flag type displayed in the help message.
	DisableUnquoteUsage bool                // DisableUnquoteUsage will toggle extract and unquote the type from the usage.
	DisablePrintDefault bool                // DisablePrintDefault toggles printing of the default value in usage message.
	Value               Value               // Value of the value as set.
	AddNegative         bool                // AddNegative automatically add a --no-<flag> option for boolean flags.
	DefValue            string              // DefValue should contain the default value (as text); for usage message.
	Changed             bool                // Changed contains whether the user set the value (or if left to default).
	Deprecated          string              // Deprecated is a string printed for a deprecation notice.
	Hidden              bool                // Hidden is used by zulu.Command to allow flags to be hidden from help/usage text.
	Required            bool                // Required ensures that a flag must be changed.
	ShorthandDeprecated string              // ShorthandDeprecated is a string printed for a deprecation notice of the Shorthand.
	Group               string              // Group contains the flag group.
	Annotations         map[string][]string // Annotations are used to annotate this specific flag for your application; e.g. it is used by zulu.Command bash completion code.
}

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
type Value interface {
	String() string
	Set(string) error
}

type Getter interface {
	Value
	Get() interface{}
}

// Typed is an interface of Values that can communicate their type.
type Typed interface {
	Type() string
}

// SliceValue is a secondary interface to all flags which hold a list
// of values.  This allows full control over the value of list flags,
// and avoids complicated marshalling and unmarshalling to csv.
type SliceValue interface {
	// Append adds the specified value to the end of the flag value list.
	Append(string) error
	// Replace will fully overwrite any data currently in the flag value list.
	Replace([]string) error
	// GetSlice returns the flag value list as an array of strings.
	GetSlice() []string
}

// MapValue is a secondary interface to all flags which hold a map of values.
type MapValue interface {
	Value
	IsMap() bool
}

// BoolFlag is an optional interface to indicate boolean flags that can be
// supplied without a value text
type BoolFlag interface {
	Value
	IsBoolFlag() bool
}

type OptionalValue interface {
	Value
	IsOptional() bool
}

// sortFlags returns the flags as a slice in lexicographical sorted order.
func sortFlags(flags map[NormalizedName]*Flag) []*Flag {
	list := make(sort.StringSlice, len(flags))
	i := 0
	for k := range flags {
		list[i] = string(k)
		i++
	}
	list.Sort()
	result := make([]*Flag, len(list))
	for i, name := range list {
		result[i] = flags[NormalizedName(name)]
	}
	return result
}

// SetNormalizeFunc allows you to add a function which can translate flag names.
// Flags added to the FlagSet will be translated and then when anything tries to
// look up the flag that will also be translated. So it would be possible to create
// a flag named "getURL" and have it translated to "geturl".  A user could then pass
// "--getUrl" which may also be translated to "geturl" and everything will work.
func (fs *FlagSet) SetNormalizeFunc(n func(f *FlagSet, name string) NormalizedName) {
	fs.normalizeNameFunc = n
	fs.sortedFormal = fs.sortedFormal[:0]
	for fname, flag := range fs.formal {
		nname := fs.normalizeFlagName(flag.Name)
		if fname == nname {
			continue
		}
		flag.Name = string(nname)
		delete(fs.formal, fname)
		fs.formal[nname] = flag
		if _, set := fs.actual[fname]; set {
			delete(fs.actual, fname)
			fs.actual[nname] = flag
		}
	}
}

// GetNormalizeFunc returns the previously set NormalizeFunc of a function which
// does no translation, if not set previously.
func (fs *FlagSet) GetNormalizeFunc() func(f *FlagSet, name string) NormalizedName {
	if fs.normalizeNameFunc != nil {
		return fs.normalizeNameFunc
	}
	return func(_ *FlagSet, name string) NormalizedName { return NormalizedName(name) }
}

func (fs *FlagSet) normalizeFlagName(name string) NormalizedName {
	n := fs.GetNormalizeFunc()
	return n(fs, name)
}

// Output returns the destination for usage and error messages. os.Stderr is returned if
// output was not set or was set to nil.
func (fs *FlagSet) Output() io.Writer {
	if fs.output == nil {
		return os.Stderr
	}
	return fs.output
}

// Name returns the name of the flag set.
func (fs *FlagSet) Name() string {
	return fs.name
}

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (fs *FlagSet) SetOutput(output io.Writer) {
	fs.output = output
}

// GetAllFlags return the flags in lexicographical order or
// in primordial order if f.SortFlags is false.
// It visits all flags, even those not set.
func (fs *FlagSet) GetAllFlags() (flags []*Flag) {
	if fs.SortFlags {
		if len(fs.formal) != len(fs.sortedFormal) {
			fs.sortedFormal = sortFlags(fs.formal)
		}
		flags = fs.sortedFormal
	} else {
		flags = fs.orderedFormal
	}
	return
}

// VisitAll visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func (fs *FlagSet) VisitAll(fn func(*Flag)) {
	if len(fs.formal) == 0 {
		return
	}
	for _, flag := range fs.GetAllFlags() {
		fn(flag)
	}
}

// HasFlags returns a bool to indicate if the FlagSet has any flags defined.
func (fs *FlagSet) HasFlags() bool {
	return len(fs.formal) > 0
}

// HasAvailableFlags returns a bool to indicate if the FlagSet has any flags
// that are not hidden.
func (fs *FlagSet) HasAvailableFlags() bool {
	for _, flag := range fs.formal {
		if !flag.Hidden {
			return true
		}
	}
	return false
}

// GetAllFlags return the flags in lexicographical order or
// in primordial order if f.SortFlags is false.
func GetAllFlags() []*Flag {
	return CommandLine.GetAllFlags()
}

// VisitAll visits the command-line flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits all flags, even those not set.
func VisitAll(fn func(*Flag)) {
	CommandLine.VisitAll(fn)
}

// GetFlags return the flags in lexicographical order or
// in primordial order if f.SortFlags is false.
// It visits only those flags that have been set.
func (fs *FlagSet) GetFlags() (flags []*Flag) {
	if fs.SortFlags {
		if len(fs.actual) != len(fs.sortedActual) {
			fs.sortedActual = sortFlags(fs.actual)
		}
		flags = fs.sortedActual
	} else {
		flags = fs.orderedActual
	}
	return
}

// Visit visits the flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func (fs *FlagSet) Visit(fn func(*Flag)) {
	if len(fs.actual) == 0 {
		return
	}
	for _, flag := range fs.GetFlags() {
		fn(flag)
	}
}

// GetFlags return the flags in lexicographical order or
// in primordial order if f.SortFlags is false.
func GetFlags() []*Flag {
	return CommandLine.GetFlags()
}

// Visit visits the command-line flags in lexicographical order or
// in primordial order if f.SortFlags is false, calling fn for each.
// It visits only those flags that have been set.
func Visit(fn func(*Flag)) {
	CommandLine.Visit(fn)
}

func (fs *FlagSet) addUnknownFlag(s string) {
	fs.unknownFlags = append(fs.unknownFlags, s)
}

// GetUnknownFlags returns unknown flags in the order they were Parsed.
// This requires ParseErrorsWhitelist.UnknownFlags to be set so that parsing does
// not abort on the first unknown flag.
func (fs *FlagSet) GetUnknownFlags() []string {
	return fs.unknownFlags
}

// GetUnknownFlags returns unknown command-line flags in the order they were Parsed.
// This requires ParseErrorsWhitelist.UnknownFlags to be set so that parsing does
// not abort on the first unknown flag.
func GetUnknownFlags() []string {
	return CommandLine.GetUnknownFlags()
}

// Get returns the value of the named flag.
func (fs *FlagSet) Get(name string) (interface{}, error) {
	return fs.getFlagValue(name, "")
}

// Get returns the value of the named flag.
func Get(name string) (interface{}, error) {
	return CommandLine.Get(name)
}

// Lookup returns the Flag structure of the named flag, returning nil if none exists.
func (fs *FlagSet) Lookup(name string) *Flag {
	return fs.lookup(fs.normalizeFlagName(name))
}

// ShorthandLookup returns the Flag structure of the shorthand flag,
// returning nil if none exists.
func (fs *FlagSet) ShorthandLookup(name rune) *Flag {
	if name == 0 {
		return nil
	}

	v, ok := fs.shorthands[name]
	if !ok {
		return nil
	}
	return v
}

// ShorthandLookupStr is the same as ShorthandLookup, but you can look it up through a string.
// It panics if name contains more than one UTF-8 character.
func (fs *FlagSet) ShorthandLookupStr(name string) *Flag {
	r, err := shorthandStrToRune(name)
	if err != nil {
		fmt.Fprintln(fs.Output(), err)
		panic(err)
	}

	return fs.ShorthandLookup(r)
}

func shorthandStrToRune(name string) (rune, error) {
	if utf8.RuneCountInString(name) > 1 {
		return 0, fmt.Errorf("cannot convert shorthand with more than one UTF-8 character: %q", name)
	}
	r, _ := utf8.DecodeRuneInString(name)
	if r == utf8.RuneError {
		return 0, nil
	}

	return r, nil
}

// lookup returns the Flag structure of the named flag, returning nil if none exists.
func (fs *FlagSet) lookup(name NormalizedName) *Flag {
	return fs.formal[name]
}

// getFlagValue returns the value of a flag based on the requested name and type.
func (fs *FlagSet) getFlagValue(name string, fType string) (interface{}, error) {
	flag := fs.Lookup(name)
	if flag == nil {
		return nil, NewUnknownFlagError(name)
	}

	if v, isTyped := flag.Value.(Typed); isTyped && fType != "" && v.Type() != fType {
		return nil, fmt.Errorf("trying to get %q value of flag of type %q", fType, v.Type())
	}

	getter, ok := flag.Value.(Getter)
	if !ok {
		return nil, fmt.Errorf("flag %q does not implement the Getter interface", name)
	}

	return getter.Get(), nil
}

// ArgsLenAtDash will return the length of f.Args at the moment when a -- was
// found during arg parsing. This allows your program to know which args were
// before the -- and which came after.
func (fs *FlagSet) ArgsLenAtDash() int {
	return fs.argsLenAtDash
}

// Lookup returns the Flag structure of the named command-line flag,
// returning nil if none exists.
func Lookup(name string) *Flag {
	return CommandLine.Lookup(name)
}

// ShorthandLookup returns the Flag structure of the shorthand flag,
// returning nil if none exists.
func ShorthandLookup(name rune) *Flag {
	return CommandLine.ShorthandLookup(name)
}

// ShorthandLookupStr is the same as ShorthandLookup, but you can look it up through a string.
// It panics if name contains more than one UTF-8 character.
func ShorthandLookupStr(name string) *Flag {
	return CommandLine.ShorthandLookupStr(name)
}

// Set sets the value of the named flag.
func (fs *FlagSet) Set(name, value string) error {
	normalName := fs.normalizeFlagName(name)
	flag, ok := fs.formal[normalName]
	if !ok {
		return NewUnknownFlagError(name)
	}

	err := flag.Value.Set(value)
	if err != nil {
		return NewInvalidArgumentError(err, flag, value)
	}

	if !flag.Changed {
		if fs.actual == nil {
			fs.actual = make(map[NormalizedName]*Flag)
		}
		fs.actual[normalName] = flag
		fs.orderedActual = append(fs.orderedActual, flag)

		flag.Changed = true
	}

	if flag.Deprecated != "" {
		fmt.Fprintf(fs.Output(), "Flag --%s has been deprecated, %s\n", flag.Name, flag.Deprecated)
	}
	return nil
}

// SetAnnotation allows one to set arbitrary annotations on this flag.
// This is sometimes used by zulucmd/zulu programs which want to generate additional
// bash completion information.
func (f *Flag) SetAnnotation(key string, values []string) {
	if f.Annotations == nil {
		f.Annotations = map[string][]string{}
	}

	f.Annotations[key] = values
}

// Changed returns true if the flag was explicitly set during Parse() and false
// otherwise
func (fs *FlagSet) Changed(name string) bool {
	flag := fs.Lookup(name)
	// If a flag doesn't exist, it wasn't changed....
	if flag == nil {
		return false
	}
	return flag.Changed
}

// Set sets the value of the named command-line flag.
func Set(name, value string) error {
	return CommandLine.Set(name, value)
}

// PrintDefaults prints to standard error unless configured otherwise, the
// default values of all defined command-line flags in the set. See the
// documentation for the global function PrintDefaults for more information.
func (fs *FlagSet) PrintDefaults() {
	usages := fs.FlagUsages()
	fmt.Fprint(fs.Output(), usages)
}

// DefaultIsZeroValue returns true if the default value for this flag represents
// a zero value.
func (f *Flag) DefaultIsZeroValue() bool {
	switch f.Value.(type) {
	case BoolFlag:
		return f.DefValue == "false"
	case SliceValue:
		return f.DefValue == "[]"
	case *durationValue:
		return f.DefValue == "0s"
	case *intValue, *int8Value, *int32Value, *int64Value, *uintValue, *uint8Value, *uint16Value, *uint32Value, *uint64Value, *countValue, *float32Value, *float64Value:
		return f.DefValue == "0"
	case *stringValue:
		return f.DefValue == ""
	case *ipValue, *ipMaskValue, *ipNetValue:
		return f.DefValue == "<nil>"
	default:
		switch f.DefValue {
		case "false", "<nil>", "", "0":
			return true
		}
		return false
	}
}

// UnquoteUsage extracts a back-quoted name from the usage
// string for a flag and returns it and the un-quoted usage.
// Given "a `name` to show" it returns ("name", "a name to show").
// If there are no back quotes, the name is an educated guess of the
// type of the flag's value, or the empty string if the flag is boolean.
func UnquoteUsage(flag *Flag) (name string, usage string) {
	name = flag.UsageType
	usage = flag.Usage

	// Look for a back-quoted name, but avoid the strings package.
	if !flag.DisableUnquoteUsage {
		name, usage = unquoteBacktickFromUsage(name, usage)
	}

	if name == "" {
		name = "value" // compatibility layer to be a drop-in replacement
		if v, ok := flag.Value.(Typed); ok {
			name = v.Type()
			switch name {
			case "bool":
				name = ""
			case "boolSlice":
				name = "bools"
			case "complex128":
				name = "complex"
			case "complex128Slice":
				name = "complexes"
			case "durationSlice":
				name = "durations"
			case "float32", "float64":
				name = "float"
			case "floatSlice", "float32Slice", "float64Slice":
				name = "floats"
			case "int8", "int16", "int32", "int64":
				name = "int"
			case "intSlice", "int8Slice", "int16Slice", "int32Slice", "int64Slice":
				name = "ints"
			case "stringSlice":
				name = "strings"
			case "uint8", "uint16", "uint32", "uint64":
				name = "uint"
			case "uintSlice", "uint8Slice", "uint16Slice", "uint32Slice", "uint64Slice":
				name = "uints"
			}
		}
	}

	return
}

func unquoteBacktickFromUsage(name string, usage string) (string, string) {
	start := strings.IndexByte(usage, '`')
	if start == -1 {
		return name, usage
	}

	end := strings.IndexByte(usage[start+1:], '`')
	if end == -1 {
		return name, usage
	}
	end += start + 1 // to skip the backtick

	extracted := usage[start+1 : end]
	if name == "" {
		name = extracted
	}
	usage = usage[:start] + extracted + usage[end+1:]

	return name, usage
}

// Splits the string `s` on whitespace into an initial substring up to
// `i` runes in length and the remainder. Will go `slop` over `i` if
// that encompasses the entire string (which allows the caller to
// avoid short orphan words on the final line). If the next word is
// wider than `i`, it is returned on its own and wrapping continues
// after it.
func wrapN(i, slop int, s string) (string, string) {
	if i+slop > len(s) {
		return s, ""
	}

	w := strings.LastIndexAny(s[:i], " \t\n")
	if w <= 0 {
		next := strings.IndexAny(s, " \t\n")
		if next <= 0 {
			return s, ""
		}
		return s[:next], s[next+1:]
	}
	nlPos := strings.LastIndex(s[:i], "\n")
	if nlPos > 0 && nlPos < w {
		return s[:nlPos], s[nlPos+1:]
	}
	return s[:w], s[w+1:]
}

// Wraps the string `s` to a maximum width `w` with leading indent
// `i`. The first line is not indented (this is assumed to be done by
// caller). Pass `w` == 0 to do no wrapping
func wrap(i, w int, s string) string {
	if w == 0 {
		return strings.ReplaceAll(s, "\n", "\n"+strings.Repeat(" ", i))
	}

	// space between indent i and end of line width w into which
	// we should wrap the text.
	wrap := w - i

	var r, l string

	// Not enough space for sensible wrapping. Wrap as a block on
	// the next line instead.
	if wrap < 24 {
		i = 16
		wrap = w - i
		r += "\n" + strings.Repeat(" ", i)
	}
	// If still not enough space then don't even try to wrap.
	if wrap < 24 {
		return strings.ReplaceAll(s, "\n", r)
	}

	// Try to avoid short orphan words on the final line, by
	// allowing wrapN to go a bit over if that would fit in the
	// remainder of the line.
	slop := 5
	wrap -= slop

	// Handle first line, which is indented by the caller (or the
	// special case above)
	l, s = wrapN(wrap, slop, s)
	r += strings.ReplaceAll(l, "\n", "\n"+strings.Repeat(" ", i))

	// Now wrap the rest
	for s != "" {
		var t string

		t, s = wrapN(wrap, slop, s)
		r += "\n" + strings.Repeat(" ", i) + strings.ReplaceAll(t, "\n", "\n"+strings.Repeat(" ", i))
	}

	return r
}

func (fs *FlagSet) flagUsageFormatter() FlagUsageFormatter {
	if fs.FlagUsageFormatter != nil {
		return fs.FlagUsageFormatter
	}

	return defaultUsageFormatter
}

// FlagUsagesWrapped returns a string containing the usage information
// for all flags in the FlagSet. Wrapped to `cols` columns (0 for no
// wrapping)
func (fs *FlagSet) FlagUsagesWrapped(cols int) string {
	return fs.FlagUsagesForGroupWrapped("", cols)
}

// FlagUsagesForGroupWrapped returns a string containing the usage information
// for all flags in the FlagSet for a group. Wrapped to `cols` columns (0 for no
// wrapping).
func (fs *FlagSet) FlagUsagesForGroupWrapped(group string, cols int) string {
	usageFormatter := fs.flagUsageFormatter()

	var (
		max, maxlen int
		lines       = make(map[string][]string)
	)
	fs.VisitAll(func(flag *Flag) {
		if flag.Hidden {
			return
		}

		line, right := usageFormatter(flag)

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		line += "\x00"
		if len(line) > maxlen {
			maxlen = len(line)
		}

		line += right

		groupName := flag.Group
		if _, ok := lines[groupName]; !ok {
			lines[groupName] = make([]string, 0)
		}
		lines[groupName] = append(lines[groupName], line)
		max += len(line) + len(groupName)
	})

	buf := new(bytes.Buffer)
	buf.Grow(max)
	for _, line := range lines[group] {
		sidx := strings.Index(line, "\x00")
		spacing := strings.Repeat(" ", maxlen-sidx)
		// maxlen + 2 comes from + 1 for the \x00 and + 1 for the (deliberate) off-by-one in maxlen-sidx
		fmt.Fprintln(buf, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
	}

	return buf.String()
}

// FlagUsages returns a string containing the usage information for all flags in
// the FlagSet
func (fs *FlagSet) FlagUsages() string {
	return fs.FlagUsagesWrapped(0)
}

// FlagUsagesForGroup returns a string containing the usage information for all flags in
// the FlagSet for group
func (fs *FlagSet) FlagUsagesForGroup(group string) string {
	return fs.FlagUsagesForGroupWrapped(group, 0)
}

// Groups return an array of unique flag groups sorted in the same order
// as flags. Empty group (unassigned) is always placed at the beginning.
func (fs *FlagSet) Groups() []string {
	groupsMap := make(map[string]bool)
	groups := make([]string, 0)
	hasUngrouped := false
	fs.VisitAll(func(flag *Flag) {
		if flag.Group == "" {
			hasUngrouped = true
			return
		}
		if _, ok := groupsMap[flag.Group]; !ok {
			groupsMap[flag.Group] = true
			groups = append(groups, flag.Group)
		}
	})
	sort.Strings(groups)

	if hasUngrouped {
		groups = append([]string{""}, groups...)
	}

	return groups
}

// PrintDefaults prints, to standard error unless configured otherwise,
// a usage message showing the default settings of all defined
// command-line flags.
// For an integer valued flag x, the default output has the form
//
//	-x int
//		usage-message-for-x (default 7)
//
// The usage message will appear on a separate line for anything but
// a bool flag with a one-byte name. For bool flags, the type is
// omitted and if the flag name is one byte the usage message appears
// on the same line. The parenthetical default is omitted if the
// default is the zero value for the type. The listed type, here int,
// can be changed by placing a back-quoted name in the flag's usage
// string; the first such item in the message is taken to be a parameter
// name to show in the message and the back quotes are stripped from
// the message when displayed. For instance, given
//
//	flag.String("I", "", "search `directory` for include files")
//
// the output will be
//
//	-I directory
//		search directory for include files.
//
// To change the destination for flag messages, call CommandLine.SetOutput.
func PrintDefaults() {
	CommandLine.PrintDefaults()
}

// defaultUsage is the default function to print a usage message.
func (fs *FlagSet) defaultUsage() {
	if fs.name == "" {
		fmt.Fprintf(fs.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.name)
	}
	fs.PrintDefaults()
}

// NOTE: Usage is not just CommandLine.defaultUsage()
// because it serves (via godoc flag Usage) as the example
// for how to write your own usage function.

// Usage prints to standard error a usage message documenting all defined command-line flags.
// The function is a variable that may be changed to point to a custom function.
// By default it prints a simple header and calls PrintDefaults; for details about the
// format of the output and how to control it, see the documentation for PrintDefaults.
var Usage = func() {
	fmt.Fprintf(CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}

// NFlag returns the number of flags that have been set.
func (fs *FlagSet) NFlag() int { return len(fs.actual) }

// NFlag returns the number of command-line flags that have been set.
func NFlag() int { return len(CommandLine.actual) }

// Arg returns the i'th argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func (fs *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(fs.args) {
		return ""
	}
	return fs.args[i]
}

// Arg returns the i'th command-line argument.  Arg(0) is the first remaining argument
// after flags have been processed.
func Arg(i int) string {
	return CommandLine.Arg(i)
}

// NArg is the number of arguments remaining after flags have been processed.
func (fs *FlagSet) NArg() int { return len(fs.args) }

// NArg is the number of arguments remaining after flags have been processed.
func NArg() int { return len(CommandLine.args) }

// Args returns the non-flag arguments.
func (fs *FlagSet) Args() []string { return fs.args }

// Args returns the non-flag command-line arguments.
func Args() []string { return CommandLine.args }

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a string into a slice of strings by
// giving the slice the methods of Value; in particular, Set would decompose
// the string into the slice.
func (fs *FlagSet) Var(value Value, name, usage string, opts ...Opt) *Flag {
	flag := &Flag{
		Name:     name,
		Usage:    usage,
		Value:    value,
		DefValue: value.String(),
	}

	if err := applyFlagOptions(flag, opts...); err != nil {
		panic(err)
	}

	fs.AddFlag(flag)
	return flag
}

// AddFlag will add the flag to the FlagSet
func (fs *FlagSet) AddFlag(flag *Flag) {
	normalizedFlagName := fs.normalizeFlagName(flag.Name)

	_, alreadyThere := fs.formal[normalizedFlagName]
	if alreadyThere {
		msg := fmt.Sprintf("%s flag redefined: %s", fs.name, flag.Name)
		fmt.Fprintln(fs.Output(), msg)
		panic(msg) // Happens only if flags are declared with identical names
	}
	if fs.formal == nil {
		fs.formal = make(map[NormalizedName]*Flag)
	}

	flag.Name = string(normalizedFlagName)
	fs.formal[normalizedFlagName] = flag
	fs.orderedFormal = append(fs.orderedFormal, flag)

	if flag.Shorthand == 0 {
		return
	}
	if fs.shorthands == nil {
		fs.shorthands = make(map[rune]*Flag)
	}
	used, alreadyThere := fs.shorthands[flag.Shorthand]
	if alreadyThere {
		msg := fmt.Sprintf("unable to redefine %q shorthand in %q flagset: it's already used for %q flag", flag.Shorthand, fs.name, used.Name)
		fmt.Fprintln(fs.Output(), msg)
		panic(msg)
	}
	fs.shorthands[flag.Shorthand] = flag
}

// RemoveFlag will remove the flag from the FlagSet
func (fs *FlagSet) RemoveFlag(name string) {
	normalizedFlagName := fs.normalizeFlagName(name)
	_, exists := fs.formal[normalizedFlagName]
	if exists {
		delete(fs.formal, normalizedFlagName)
	}
}

// AddFlagSet adds one FlagSet to another. If a flag is already present in f
// the flag from newSet will be ignored.
func (fs *FlagSet) AddFlagSet(newSet *FlagSet) {
	if newSet == nil {
		return
	}
	newSet.VisitAll(func(flag *Flag) {
		if fs.Lookup(flag.Name) == nil {
			fs.AddFlag(flag)
		}
	})
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a string into a slice of strings by
// giving the slice the methods of Value; in particular, Set would decompose
// the string into the slice.
func Var(value Value, name, usage string, opts ...Opt) *Flag {
	return CommandLine.Var(value, name, usage, opts...)
}

// failf prints to standard error a formatted error and usage message and
// returns the error.
func (fs *FlagSet) failf(format string, a ...interface{}) error {
	fs.usage()
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(fs.Output())
	fmt.Fprintln(fs.Output(), err)
	return err
}

// usage calls the Usage method for the flag set, or the usage function if
// the flag set is CommandLine.
func (fs *FlagSet) usage() {
	switch {
	case fs == CommandLine:
		Usage()
	case fs.Usage == nil:
		fs.defaultUsage()
	default:
		fs.Usage()
	}
}

// --unknown (args will be empty)
// --unknown --next-flag ... (args will be --next-flag ...)
// --unknown arg ... (args will be arg ...)
func (fs *FlagSet) stripUnknownFlagValue(args []string) []string {
	if len(args) == 0 {
		// --unknown
		return args
	}

	first := args[0]
	if len(first) > 0 && first[0] == '-' {
		// --unknown --next-flag ...
		return args
	}

	// --unknown arg ... (args will be arg ...)
	if len(args) > 1 {
		fs.addUnknownFlag(args[0])
		return args[1:]
	}
	return nil
}

//nolint:funlen
func (fs *FlagSet) parseLongArg(s string, args []string, fn parseFunc) (outArgs []string, err error) {
	outArgs = args
	name := s[2:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		err = fs.failf("bad flag syntax: %s", s)
		return
	}

	hasNoPrefix := strings.HasPrefix(name, "no-")
	split := strings.SplitN(name, "=", 2)
	name = split[0]
	flag, exists := fs.formal[fs.normalizeFlagName(name)]

	if !exists && len(name) > 3 && hasNoPrefix {
		bFlag, bExists := fs.formal[fs.normalizeFlagName(name[3:])]
		if bExists && bFlag.AddNegative {
			if _, isBoolFlag := bFlag.Value.(BoolFlag); isBoolFlag {
				flag = bFlag
				exists = bExists
				name = name[3:]
			}
		}
	}

	if !exists || (flag != nil && flag.ShorthandOnly) {
		switch {
		case !exists && name == "help" && !fs.DisableBuiltinHelp:
			fs.usage()
			err = ErrHelp
			return
		case fs.ParseErrorsAllowList.UnknownFlags || (flag != nil && flag.ShorthandOnly):
			// --unknown=unknownval arg ...
			// we do not want to lose arg in this case
			fs.addUnknownFlag(s)
			if len(split) >= 2 {
				return
			}
			outArgs = fs.stripUnknownFlagValue(outArgs)
			return
		default:
			err = fs.failf("%s", NewUnknownFlagError(name).Error())
			return
		}
	}

	_, flagIsBool := flag.Value.(BoolFlag)
	_, isOptional := flag.Value.(OptionalValue)
	nextArgIsFlagValue := len(outArgs) > 0 && len(outArgs[0]) > 0 && outArgs[0][0] != '-'

	var value string
	switch {
	case len(split) == 2: // '--flag=arg'
		value = split[1]
		if hasNoPrefix && flagIsBool {
			err = fs.failf("flag cannot have a value: %s", s)
			return
		}
	case flagIsBool: // '--[no-]flag' (arg was optional)
		value = fmt.Sprintf("%t", !hasNoPrefix)
	case isOptional: // '--flag' (arg was optional)
		value = ""
	case nextArgIsFlagValue && (!flagIsBool || (flagIsBool && isBool(outArgs[0]))): // '--flag arg'
		value = outArgs[0]
		outArgs = outArgs[1:]
	default: // '--flag' (arg was required)
		err = fs.failf("flag needs an argument: %s", s)
		return
	}

	err = fn(flag, value)
	if err != nil {
		err = fs.failf("%s", err.Error())
	}
	return
}

func isBool(v string) bool {
	_, err := strconv.ParseBool(v)
	return err == nil
}

//nolint:funlen
func (fs *FlagSet) parseSingleShortArg(shorthands string, args []string, fn parseFunc) (outShorts string, outArgs []string, err error) {
	outArgs = args
	outShorts = shorthands[1:]
	char, _ := utf8.DecodeRuneInString(shorthands)

	flag, exists := fs.shorthands[char]
	if !exists {
		switch {
		case char == 'h' && !fs.DisableBuiltinHelp:
			fs.usage()
			err = ErrHelp
			return
		case fs.ParseErrorsAllowList.UnknownFlags:
			if len(shorthands) > 2 {
				// '-f...'
				// we do not want to lose anything in this case
				fs.addUnknownFlag("-" + shorthands)
				outShorts = ""
				return
			}
			fs.addUnknownFlag("-" + string(char))
			if len(outShorts) == 0 {
				outArgs = fs.stripUnknownFlagValue(outArgs)
			}
			return
		default:
			// fallback to a normal flag look up without any shorthand opts
			flag = fs.Lookup(string(char))
			if flag == nil || (flag.Shorthand > 0 && flag.Shorthand != char) {
				err = fs.failf("unknown shorthand flag: %q in -%s", char, shorthands)
				return
			}
		}
	}

	_, flagIsBool := flag.Value.(BoolFlag)
	_, isOptional := flag.Value.(OptionalValue)
	nextArgIsFlagValue := len(outArgs) > 0 && len(outArgs[0]) > 0 && outArgs[0][0] != '-'

	nextShortArgIsFlagValue := len(shorthands) > 1
	if len(shorthands) > 1 {
		_, nextFlagExists := fs.shorthands[rune(shorthands[1])]
		nextShortArgIsFlagValue = !nextFlagExists
	}

	var value string
	switch {
	case len(shorthands) > 2 && shorthands[1] == '=':
		// '-f=arg'
		value = shorthands[2:]
		outShorts = ""
	case nextShortArgIsFlagValue && (!flagIsBool || (flagIsBool && isBool(shorthands[1:]))):
		// '-farg'
		value = shorthands[1:]
		outShorts = ""
	case nextArgIsFlagValue && (!flagIsBool || (flagIsBool && isBool(outArgs[0]))):
		// '-f arg'
		value = args[0]
		outArgs = args[1:]
	case flagIsBool, isOptional:
		// '-f' (arg was optional)
		value = ""
	default:
		// '-f' (arg was required)
		err = fs.failf("flag needs an argument: %q in -%s", char, shorthands)
		return
	}

	if flag.ShorthandDeprecated != "" {
		fmt.Fprintf(fs.Output(), "Flag shorthand -%c has been deprecated, %s\n", flag.Shorthand, flag.ShorthandDeprecated)
	}

	err = fn(flag, value)
	if err != nil {
		err = fs.failf("%s", err.Error())
	}
	return
}

func (fs *FlagSet) parseShortArg(s string, args []string, fn parseFunc) (outArgs []string, err error) {
	outArgs = args
	shorthands := s[1:]

	// "shorthands" can be a series of shorthand letters of flags (e.g. "-vvv").
	for utf8.RuneCountInString(shorthands) > 0 {
		shorthands, outArgs, err = fs.parseSingleShortArg(shorthands, args, fn)
		if err != nil {
			return
		}
	}

	return
}

func (fs *FlagSet) parseArgs(args []string, fn parseFunc) (err error) {
	for len(args) > 0 {
		s := args[0]
		args = args[1:]
		if len(s) == 0 || s[0] != '-' || len(s) == 1 {
			if !fs.interspersed {
				fs.args = append(fs.args, s)
				fs.args = append(fs.args, args...)
				return nil
			}
			fs.args = append(fs.args, s)
			continue
		}

		if s[1] == '-' {
			if len(s) == 2 && s == "--" { // "--" terminates the flags
				fs.argsLenAtDash = len(fs.args)
				fs.args = append(fs.args, args...)
				break
			}
			args, err = fs.parseLongArg(s, args, fn)
		} else {
			args, err = fs.parseShortArg(s, args, fn)
		}
		if err != nil {
			return
		}
	}

	return fs.Validate()
}

var exitFn = func(code int) {
	os.Exit(code)
}

func (fs *FlagSet) parseAll(arguments []string, fn parseFunc) error {
	if fs.addedGoFlagSets != nil {
		for _, goFlagSet := range fs.addedGoFlagSets {
			if err := goFlagSet.Parse(nil); err != nil {
				return err
			}
		}
	}
	fs.parsed = true

	fs.args = make([]string, 0, len(arguments))

	if len(arguments) == 0 {
		return fs.Validate()
	}

	err := fs.parseArgs(arguments, fn)
	if err != nil {
		switch fs.errorHandling {
		case ContinueOnError:
			return err
		case ExitOnError:
			if err == ErrHelp {
				exitFn(0)
			}
			exitFn(2)
		case PanicOnError:
			panic(err)
		}
	}
	return nil
}

// Parse parses flag definitions from the argument list, which should not
// include the command name.  Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help was set but not defined.
func (fs *FlagSet) Parse(arguments []string) error {
	set := func(flag *Flag, value string) error {
		return fs.Set(flag.Name, value)
	}
	return fs.parseAll(arguments, set)
}

type parseFunc func(flag *Flag, value string) error

// ParseAll parses flag definitions from the argument list, which should not
// include the command name. The arguments for fn are flag and value. Must be
// called after all flags in the FlagSet are defined and before flags are
// accessed by the program. The return value will be ErrHelp if -help was set
// but not defined.
func (fs *FlagSet) ParseAll(arguments []string, fn func(flag *Flag, value string) error) error {
	return fs.parseAll(arguments, fn)
}

// Parsed reports whether f.Parse has been called.
func (fs *FlagSet) Parsed() bool {
	return fs.parsed
}

// Parse parses the command-line flags from os.Args[1:].  Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() {
	// Ignore errors; CommandLine is set for ExitOnError.
	_ = CommandLine.Parse(os.Args[1:])
}

// ParseAll parses the command-line flags from os.Args[1:] and called fn for each.
// The arguments for fn are flag and value. Must be called after all flags are
// defined and before flags are accessed by the program.
func ParseAll(fn func(flag *Flag, value string) error) {
	// Ignore errors; CommandLine is set for ExitOnError.
	_ = CommandLine.ParseAll(os.Args[1:], fn)
}

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func SetInterspersed(interspersed bool) {
	CommandLine.SetInterspersed(interspersed)
}

// Parsed returns true if the command-line flags have been parsed.
func Parsed() bool {
	return CommandLine.Parsed()
}

// CommandLine is the default set of command-line flags, parsed from os.Args.
var CommandLine = NewFlagSet(os.Args[0], ExitOnError)

// NewFlagSet returns a new, empty flag set with the specified name,
// error handling property and SortFlags set to true.
func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet {
	f := &FlagSet{
		name:          name,
		errorHandling: errorHandling,
		argsLenAtDash: -1,
		interspersed:  true,
		SortFlags:     true,
	}
	return f
}

// SetInterspersed sets whether to support interspersed option/non-option arguments.
func (fs *FlagSet) SetInterspersed(interspersed bool) {
	fs.interspersed = interspersed
}

// Init sets the name and error handling property for a flag set.
// By default, the zero FlagSet uses an empty name and the
// ContinueOnError error handling policy.
func (fs *FlagSet) Init(name string, errorHandling ErrorHandling) {
	fs.name = name
	fs.errorHandling = errorHandling
	fs.argsLenAtDash = -1
}

// Validate ensures all flag values are valid.
func (fs *FlagSet) Validate() error {
	if !fs.ParseErrorsAllowList.RequiredFlags {
		var missingFlagsErr MissingFlagsError
		fs.VisitAll(func(f *Flag) {
			if f.Required && !f.Changed {
				missingFlagsErr.AddMissingFlag(f)
			}
		})

		if len(missingFlagsErr) > 0 {
			return missingFlagsErr
		}
	}

	return nil
}
