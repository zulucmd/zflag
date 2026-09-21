# zflag

[![GoDoc](https://godoc.org/github.com/zulucmd/zflag?status.svg)](https://godoc.org/github.com/zulucmd/zflag)
[![Go Report Card](https://goreportcard.com/badge/github.com/zulucmd/zflag)](https://goreportcard.com/report/github.com/zulucmd/zflag)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/zulucmd/zflag?sort=semver)](https://github.com/zulucmd/zflag/releases)
[![Build Status](https://github.com/zulucmd/zflag/actions/workflows/validate.yml/badge.svg)](https://github.com/zulucmd/zflag/actions/workflows/validate.yml)

<!-- toc -->
- [Installation](#installation)
- [Supported Syntax](#supported-syntax)
- [Fork from pflag](#fork-from-pflag)
- [Documentation](#documentation)
  - [Quick start](#quick-start)
  - [Bool Values](#bool-values)
  - [Mutating or &quot;Normalizing&quot; Flag names](#mutating-or-normalizing-flag-names)
  - [Deprecating a flag or its shorthand](#deprecating-a-flag-or-its-shorthand)
  - [Hidden flags](#hidden-flags)
  - [Required flags](#required-flags)
  - [Disable sorting of flags](#disable-sorting-of-flags)
  - [Supporting Go flags when using zflag](#supporting-go-flags-when-using-zflag)
  - [Shorthand flags](#shorthand-flags)
  - [Shorthand-only flags](#shorthand-only-flags)
  - [Unknown flags](#unknown-flags)
  - [Custom flag types in usage](#custom-flag-types-in-usage)
  - [Customizing flag usages](#customizing-flag-usages)
  - [Disable printing a flag's default value](#disable-printing-a-flags-default-value)
  - [Disable built-in help flags](#disable-built-in-help-flags)
- [Development](#development)
<!-- /toc -->

## Installation

zflag is available using the standard `go get` command.

Install by running:

```bash
go get github.com/zulucmd/zflag/v2
```

## Supported Syntax

```plain
--flag       // boolean flags, or flags with no option default values
--no-flag    // boolean flags
--flag x
--flag=x
```

Unlike the flag package, a single dash before an option means something
different than a double dash. Single dashes signify a series of shorthand
letters for flags. All but the last shorthand letter must be boolean flags
or a flag with a default value

```plain
// boolean or flags where the 'no option default value' is set
-f
-f=true
-f true
-abc

// non-boolean and flags without a 'no option default value'
-n 1234
-n=1234
-n1234

// mixed
-abcs "hello"
-absd="hello"
-abcs1234
```

Slice flags can be specified multiple times.

```plain
--sliceVal one --sliceVal=two
```

Mapped flags can be specified.

```plain
--map-val key1=value --map-val key2=value
```

Integer flags accept 1234, 0664, 0x1234 and may be negative.
Boolean flags accept 1, 0, t, f, true, false,
TRUE, FALSE, True, False.
Duration flags accept any input valid for time.ParseDuration.

Flag parsing stops after the terminator "--". Unlike the flag package,
flags can be interspersed with arguments anywhere on the command line
before this terminator.

## Fork from pflag

This is a fork of [cornfeedhobo/pflag](https://github.com/cornfeedhobo/pflag), which in turn is a fork of [spf13/pflag](https://github.com/spf13/pflag).

Both repos haven't had any updates or maintenance. I'm sure for many that's fine, but I needed changes
that weren't available.

The following are differences between zflag and pflag:

A bunch of PRs have been merged:
- Add support for flag groups.
- Switched from "whitelist" terminology to "allowlist".
- Improve flag errors dashes.
- Move Type() into its own interface.
- Use `DefValue` for usage.
- Store shorthand flag as a single UTF-8 character (`rune`).
- Refactored flag parsing to be based on a Getter interface.

In addition to the above PRs, the following changes have been made:
- A new flag usage formatter has been added to customize the usage output.
- Unknown flag errors are now consistent.
- Removed all the CSV parsing in slice types and others. These were causing more head-ache than needed,
  as it is hard to get this right for a wide variety of use cases. If you need this, please use either
  use the `Func` flag type, or creating your own custom flag type.
- Improved go `flag` compatibility:
  - Standardized the flag API. This follows the `flag` closer. Additional options can be added using `Opt*` method calls.
  - Added a `Func` flag type.
  - The `Getter` interface was implemented.
- Restructured the tests to be based on test tables and enable `t.Parallel()` where possible.
- Removed the `NoOptDefVal` in favour of interfaces.

## Documentation

You can see the full reference documentation of the zflag package
[at godoc.org](http://godoc.org/github.com/zulucmd/zflag), querying with
[`go doc`](https://golang.org/cmd/doc/), or through go's standard documentation
system by running `godoc -http=:6060` and browsing to
[http://localhost:6060/pkg/github.com/zulucmd/zflag](http://localhost:6060/pkg/github.com/zulucmd/zflag)
after installation.

### Quick start
To quickly create a parser, you can simply use the package global flag set by calling
the relevant functions:
```go
zflag.Bool("mybool", false, "my usage")
zflag.Parse()

v, err := zflag.Get("mybool")
mybool := v.(bool)
// alternatively you can also interact directly with zflag.CommandLine
mybool, err := zflag.CommandLine.GetBool("mybool")

// or

var mybool bool
zflag.BoolVar(&mybool, "mybool", false, "my usage")
zflag.Parse()

// do something with mybool
```

You can also create a custom `FlagSet` instead, this allows you not rely on global states
and allows for easier writing of co-routine safe code.

```go
f := zflag.NewFlagSet("mycustomapp", zflag.ExitOnError)
f.Bool("mybool", false, "My bool")
err := f.Parse(os.Args[1:])

mybool := f.GetBool("mybool")
```

### Bool Values

If a bool flag is added, both `--flag-name` and `--no-flag-name` will be accepted.
When using `--flag-name` the value is set to true, when using `--no-flag-name` the
value is set to false.

**Example**:

```go
var enable = flag.Bool("enable", false, "help message", flag.OptShorthand('f'))
```

**Results**:

| Parsed Arguments | Resulting Value |
|------------------|-----------------|
| --enable         | enable=true     |
| --no-enable      | enable=false    |
| [nothing]        | enable=false    |

### Mutating or "Normalizing" Flag names

It is possible to set a custom flag name 'normalization function.' It allows
flag names to be mutated both when created in the code and when used on the
command line to some 'normalized' form. The 'normalized' form is used for
comparison. Two examples of using the custom normalization func follow.

**Example #1**: You want -, _, and . in flags to compare the same. aka --my-flag == --my_flag == --my.flag

```go
func wordSepNormalizeFunc(f *zflag.FlagSet, name string) zflag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return zflag.NormalizedName(name)
}

myFlagSet.SetNormalizeFunc(wordSepNormalizeFunc)
```

**Example #2**: You want to alias two flags. aka --old-flag-name == --new-flag-name

```go
func aliasNormalizeFunc(f *zflag.FlagSet, name string) zflag.NormalizedName {
	switch name {
	case "old-flag-name":
		name = "new-flag-name"
		break
	}
	return zflag.NormalizedName(name)
}

myFlagSet.SetNormalizeFunc(aliasNormalizeFunc)
```

### Deprecating a flag or its shorthand

It is possible to deprecate a flag, or just its shorthand. Deprecating a
flag/shorthand hides it from help text and prints a usage message when the
deprecated flag/shorthand is used.

**Example #1**: You want to deprecate a flag named "badflag" as well as
inform the users what flag they should use instead.

```go
// deprecate a flag by specifying its name and a usage message
flags.Bool("badflag", false, "this does something", zflag.OptDeprecated("please use --good-flag instead"))
```

This hides "badflag" from help text, and prints
`Flag --badflag has been deprecated, please use --good-flag instead`
when "badflag" is used.

**Example #2**: You want to keep a flag name "noshorthandflag" but deprecate
it's shortname "n".

```go
// deprecate a flag shorthand by specifying its flag name and a usage message
flags.Bool("noshorthandflag", false, "this does something", zflag.OptShorthand("n"), zflag.OptShorthandDeprecated("please use --noshorthandflag only"))
```

This hides the shortname "n" from help text, and prints
`Flag shorthand -n has been deprecated, please use --noshorthandflag only`
when the shorthand `n` is used.

Note that usage message is essential here, and it should not be empty. If it is empty, it will panic at runtime.

### Hidden flags

It is possible to mark a flag as hidden, meaning it will still function as
normal, however will not show up in usage/help text.

**Example**: You have a flag named "secretFlag" that you need for internal use
only and don't want it showing up in help text, or for its usage text to be available.

```go
// hide a flag by specifying its name
flags.Bool("secretFlag", false, "this does something", zflag.OptHidden())
```

### Required flags

It is possible to mark a flag as required, meaning it zflag will return an error if
it is not passed in.

**Example**:

```go
flags.Bool("must", false, "this does something", zflag.OptRequired())
err := flags.Parse()
// err == `required flag(s) "--must" not set`
```

### Disable sorting of flags

It is possible to disable sorting of flags for help and usage message.

**Example**:

```go
flag.Bool("verbose", false, "verbose output", flag.OptShorthand('v'))
flag.String("coolflag", "yeaah", "it's really cool flag")
flag.Int("usefulflag", 777, "sometimes it's very useful")
flag.SortFlags = false
flag.PrintDefaults()
```

**Output**:

```plain
  -v, --verbose           verbose output
      --coolflag string   it's really cool flag (default "yeaah")
      --usefulflag int    sometimes it's very useful (default 777)
```

### Supporting Go flags when using zflag

In order to support flags defined using Go's `flag` package, they must be added
to the `zflag` flagset. This is usually necessary to support flags defined by
third-party dependencies (e.g. `golang/glog`).

**Example**: You want to add the Go flags to the `CommandLine` flagset

```go
import (
	goflag "flag"
	flag "github.com/zulucmd/zflag/v2"
)

var ip *int = flag.Int("flagname", 1234, "help message for flagname")

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}
```

### Shorthand flags

A flag supporting both long and short formats can be created with any of the
flag functions suffixed with `P`:

```go
flag.Bool("toggle", false, "toggle help message", zflag.OptShorthand('t'))
```

### Shorthand-only flags

A shorthand-only flag can be created with any of the flag functions suffixed
with `S`:

```go
flag.String("value", "", "value help message", zflag.OptShorthandOnly('l'))
```

This flag can be looked up using it's long name, but will only be parsed when
the short form is passed.

### Unknown flags

Normally zflag will error when an unknown flag is passed, but it's also possible
to disable that using `FlagSet.ParseErrorsAllowlist.UnknownFlags`:

```go
flags.ParseErrorsAllowlist.UnknownFlags = true
flag.Parse()
```

These can then be obtained as a slice of strings using `FlagSet.GetUnknownFlags()`.

### Custom flag types in usage

There are two methods to set a custom type to be printed in the usage.

First, it's possible to set explicitly with `UsageType`:

```go
flag.String("character", "", "character name", zflag.OptUsageType("enum"))
```

Output:

```plain
  --character enum   character name (default "")
```

Alternatively, it's possible to include backticks around a single word in the
usage string, which will be extracted and printed with the usage:

```go
flag.String("character", "", "`character` name")
```

Output:

```plain
  --character character   character name (default "")
```

_Note: This unquoting behavior can be disabled with `Flag.DisableUnquoteUsage`, or `zflag.OptDisableUnquoteUsage`_.

### Customizing flag usages

You can customize the flag usages by overriding the `FlagSet.FlagUsageFormatter` field
with function that returns two strings. By default, it uses the [`defaultUsageFormatter`](./formatter.go).
The function must return two strings, one that contains the "left" side of the usage,
and one that returns the "right" side of the usage. The sides are there to calculate
how far to indent usage.

For example:
```go
flagSet := zflag.NewFlagSet("myapp", zflag.ExitOnError)
flagSet.String("hello", "", "myusage")

flagSet.FlagUsageFormatter = func (flag *zflag.Flag) (string, string) {
  return "--not-hello string", "different usage text"
}
flagSet.PrintDefaults()
```

Which will print:
```
--not-hello string   myusage
```

### Disable printing a flag's default value

The printing of a flag's default value can be suppressed with `Flag.DisablePrintDefault`.

**Example**:

```go
flag.Int("in", -1, "help message", zflag.OptDisablePrintDefault())
```

**Output**:

```plain
  --in int   help message
```

Note: if you override the usage formatter, you'll need to take the field `Flag.DisablePrintDefault` into account.

### Disable built-in help flags

Normally zflag will handle `--help` and `-h` when the flags aren't explicitly defined.

If for some reason there is a need to capture the error returned in this condition, it
is possible to disable this built-in handling.

```go
myFlagSet.DisableBuiltinHelp = true
```

## Development

Tooling is managed with [mise](https://mise.jdx.dev). Run `mise install`, then:

- `mise run setup` — install the git hooks (prek)
- `mise run lint` — run golangci-lint
- `mise run test` — run the test suite
- `mise run check` — run every prek hook against all files

Git hooks are run by [prek](https://github.com/j178/prek).
