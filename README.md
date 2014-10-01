# envtmpl

Executes text templates using data from environment variables and sends the result to STDOUT.

[![Build Status](https://travis-ci.org/williambailey/go-envtmpl.svg)](https://travis-ci.org/williambailey/go-envtmpl)

## Install

```
$ go get -u github.com/williambailey/go-envtmpl/envtmpl
$ $GOPATH/bin/envtmpl
Usage:
  envtmpl tmplDir tmplName.tmpl

Parse tmplDir/*.tmpl and renders tmplName.tmpl to
STDOUT using environment variables.

Help:
  envtmpl -h
```

## Usage: envtmpl tmplDir tmplName.tmpl

Parse **tmplDir/*.tmpl** and renders **tmplName.tmpl** to STDOUT using
environment variables.

See http://golang.org/pkg/text/template/ for template syntax.

See http://code.google.com/p/re2/wiki/Syntax for regular expression syntax.

### Exit codes

* 0 - OK.
* 1 - Usage.
* 2 - Template parse error.
* 3 - Template execution error.

## Template Functions

### base64Encode

Encodes a value to base64.

Template:

    {{ "Hello WORLD!" | base64Encode }}

Output:

    SGVsbG8gV09STEQh

### jsonEncode

Encodes a value to JSON.

Template:

    {{ "Hello\n<WORLD>!" | jsonEncode }}

Output:

    "Hello\n\u003cWORLD\u003e!"

Template:

    {{ split "Hello\n<WORLD>!" "\n" | jsonEncode }}

Output:

    [
      "Hello",
      "\u003cWORLD\u003e!"
    ]

### linePrefix

Prefix each line.

Template:

    {{ "line1\nline2\nline3" | linePrefix "- " }}

Output:

    - line1
    - line2
    - line3

### lower

Convert to lower case.

Template:

    {{ "Hello WORLD!" | lower }}

Output:

    hello world!

### regexReplace

Replace values using a regular expression.

Template:

    {{ "this is something" | regexReplace "(this) is " "[$1] was " }}

Output:

    [this] was something

### split

Split in a string substrings using another string.

Template:

    {{ range $k, $v := split "foo BAR bAz" " " }}{{ $k }}={{ $v }} {{ end }}

Output:

    0=foo 1=BAR 2=bAz 

### title

Convert to title case.

Template:

    {{ "foo BAR bAz" | title }}

Output:

    Foo BAR BAz

### trimPrefix

Remove leading prefix. If the string doesn't start with the prefix then it's unchanged.

Template:

    {{ "foo.bar" | trimPrefix "baz." }}

Output:

    foo.bar

Template:

    {{ "foo.bar" | trimPrefix "foo." }}

Output:

    bar

### trimSpace

Remove all leading and trailing white space.

Template:

    {{ " \t\n foo bar \t\n " | trimSpace }}

Output:

    foo bar

### trimSuffix

Remove trailing suffix. If the string doesn't end with the suffix then it's unchanged.

Template:

    {{ "foo.bar" | trimSuffix ".bar" }}

Output:

    foo

Template:

    {{ "foo.bar" | trimSuffix ".baz" }}

Output:

    foo.bar

### upper

Convert to upper case.

Template:

    {{ "Hello World!" | upper }}

Output:

    HELLO WORLD!

### uuid

Create a random (v4) UUID.

Template:

    {{ uuid }}

Output:

    51e8ed98-e2c9-490b-a0bd-9c287a7dcd7b

# Contributing

1. Fork the repository on GitHub
2. Create a named feature branch (i.e. `add-new-feature`)
3. Write your change
4. Submit a Pull Request

# Authors

- William Bailey - [@cowboysfromhell](https://twitter.com/cowboysfromhell) - ([mail@williambailey.org.uk](mailto:mail@williambailey.org.uk))

# License

Licensed under a [MIT license](LICENSE.txt).
