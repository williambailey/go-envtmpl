# envtmpl

Executes text templates using data from environment variables and sends the result to STDOUT.

[![Build Status](https://travis-ci.org/williambailey/go-envtmpl.svg)](https://travis-ci.org/williambailey/go-envtmpl)

## Install

```bash
go get -u github.com/williambailey/go-envtmpl/envtmpl
```

## Usage: envtmpl tmplDir tmplName.tmpl

Parse **tmplDir/*.tmpl** and renders **tmplName.tmpl** to STDOUT using
the data from environmental variables.

See http://golang.org/pkg/text/template/ for template syntax.

See http://code.google.com/p/re2/wiki/Syntax for regular expression syntax.

### Exit codes

* 0 - OK.
* 1 - Usage.
* 2 - Template parse error.
* 3 - Template execution error.

## Template Functions

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

    {{ "foo BAR bAz" | lower }}

Output:

    foo bar baz

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

### upper

Convert to upper case.

Template:

    {{ "foo BAR bAz" | upper }}

Output:

    FOO BAR BAZ

### uuid

Create a random (v4) UUID.

Template:

    {{ uuid }}

Output:

    dd60fce2-e929-485a-888e-db9801b95770

exit status 1

# Contributing

1. Fork the repository on GitHub
2. Create a named feature branch (i.e. `add-new-feature`)
3. Write your change
4. Submit a Pull Request

# Authors

- William Bailey - [@cowboysfromhell](https://twitter.com/cowboysfromhell) - ([mail@williambailey.org.uk](mailto:mail@williambailey.org.uk))

# License

Licensed under a [MIT license](LICENSE.txt).
