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

In addition to the [actions](http://golang.org/pkg/text/template/#hdr-Actions)
and [functions](http://golang.org/pkg/text/template/#hdr-Functions) provided by
the core [template engine](http://golang.org/pkg/text/template/#pkg-overview),
envtmpl provides the following functions for use in your templates:

* [base64Encode](#base64encode) - Encodes a value to base64.
* [jsonDecode](#jsondecode) - Decodes a JSON string.
* [jsonEncode](#jsonencode) - Encodes a value to JSON.
* [linePrefix](#lineprefix) - Prefix each line.
* [lower](#lower) - Convert to lower case.
* [regexReplace](#regexreplace) - Replace values using a regular expression.
* [slugify](#slugify) - Transform text to a slugified version.
* [split](#split) - Split in a string substrings using another string.
* [title](#title) - Convert to title case.
* [trimPrefix](#trimprefix) - Remove leading prefix.
* [trimSpace](#trimspace) - Remove all leading and trailing white space.
* [trimSuffix](#trimsuffix) - Remove trailing suffix.
* [upper](#upper) - Convert to upper case.
* [url](#url) - Parse a URL from a string.
* [uuid](#uuid) - Create a random (v4) UUID.
* [wordWrap](#wordwrap) - Wraps text to a given number of runes.

### base64Encode

Encodes a value to base64.

Template:

    {{ "Hello WORLD!" | base64Encode }}

Output:

    SGVsbG8gV09STEQh

### jsonDecode

Decodes a JSON string.

Template:

    {{ $j := "{\"foo\":\"bar\"}" | jsonDecode }}Foo is {{ $j.foo }}

Output:

    Foo is bar

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

### slugify

Transform text to a slugified version. You can optionally specify a unicode
normalization rule (NFC, NFD, NFKC, NFKD). Default is NFC.

Template:

     NFC: {{ "Hello áçćèńtš!" | slugify "NFC" }}
     NFD: {{ "Hello áçćèńtš!" | slugify "NFD" }}
    NFKC: {{ "Hello áçćèńtš!" | slugify "NFKC" }}
    NFKD: {{ "Hello áçćèńtš!" | slugify "NFKD" }}

Output:

     NFC: hello-áçćèńtš
     NFD: hello-accents
    NFKC: hello-áçćèńtš
    NFKD: hello-accents

Template:

    {{ "Hello W/O-R_L~D!" | slugify }}

Output:

    hello-wo-r_l~d

Template:

    {{ "Hello WORLD!" | slugify }}

Output:

    hello-world

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

Remove leading prefix. If the string doesn't start with the prefix then it's
unchanged.

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

Remove trailing suffix. If the string doesn't end with the suffix then it's
unchanged.

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

### url

Parse a URL from a string. You can optionally provide a base URL that will be
used as the context for processing.

Template:

    {{define "u"}}Url: {{ . }}
    IsAbs: {{ .IsAbs }}
    Scheme: "{{ .Scheme }}"
    Opaque: "{{ .Opaque }}"
    UserPassword: "{{ .UserPassword }}"
    Username: "{{ .Username }}"
    HasPassword: {{ .HasPassword }}
    Password: "{{ .Password }}"
    Host: "{{ .Host }}"
    Path: "{{ .Path }}"
    RawQuery: "{{ .RawQuery }}"
    Query: {{ .Query | jsonEncode }}
    Fragment: "{{ .Fragment }}"
    RequestURI: "{{ .RequestURI }}"
    {{end}}
    {{ template "u" "scheme://username:password@domain:port/path?query=string#fragment_id" | url }}
    {{ template "u" "scheme://username:@domain?a=1&b=2&a=11" | url }}
    {{ template "u" "scheme:opaque?query=string#fragment_id" | url }}
    {{ template "u" "../bar/baz" | url }}
    {{ template "u" "../bar/baz" | url "scheme://domain/foo/qux/" }}

Output:

    
    Url: scheme://username:password@domain:port/path?query=string#fragment_id
    IsAbs: true
    Scheme: "scheme"
    Opaque: ""
    UserPassword: "username:password"
    Username: "username"
    HasPassword: true
    Password: "password"
    Host: "domain:port"
    Path: "/path"
    RawQuery: "query=string"
    Query: {
      "query": [
        "string"
      ]
    }
    Fragment: "fragment_id"
    RequestURI: "/path?query=string"
    
    Url: scheme://username:@domain?a=1&b=2&a=11
    IsAbs: true
    Scheme: "scheme"
    Opaque: ""
    UserPassword: "username:"
    Username: "username"
    HasPassword: true
    Password: ""
    Host: "domain"
    Path: ""
    RawQuery: "a=1&b=2&a=11"
    Query: {
      "a": [
        "1",
        "11"
      ],
      "b": [
        "2"
      ]
    }
    Fragment: ""
    RequestURI: "/?a=1&b=2&a=11"
    
    Url: scheme:opaque?query=string#fragment_id
    IsAbs: true
    Scheme: "scheme"
    Opaque: "opaque"
    UserPassword: ""
    Username: ""
    HasPassword: false
    Password: ""
    Host: ""
    Path: ""
    RawQuery: "query=string"
    Query: {
      "query": [
        "string"
      ]
    }
    Fragment: "fragment_id"
    RequestURI: "opaque?query=string"
    
    Url: ../bar/baz
    IsAbs: false
    Scheme: ""
    Opaque: ""
    UserPassword: ""
    Username: ""
    HasPassword: false
    Password: ""
    Host: ""
    Path: "../bar/baz"
    RawQuery: ""
    Query: {}
    Fragment: ""
    RequestURI: "../bar/baz"
    
    Url: scheme://domain/foo/bar/baz
    IsAbs: true
    Scheme: "scheme"
    Opaque: ""
    UserPassword: ""
    Username: ""
    HasPassword: false
    Password: ""
    Host: "domain"
    Path: "/foo/bar/baz"
    RawQuery: ""
    Query: {}
    Fragment: ""
    RequestURI: "/foo/bar/baz"
    

### uuid

Create a random (v4) UUID.

Template:

    {{ uuid }}

Output:

    046e39e2-c050-4f3f-8adf-e6380639522f

### wordWrap

Wraps text to a given number of runes. Any existing white space is lost in the
transformation.

Template:

    {{ "The quick brown fox jumps over the lazy dog." | wordWrap 19 }}

Output:

    The quick brown
    fox jumps over the
    lazy dog.

Template:

    {{ "\t  The quick\nbrown fox jumps over the\n\t\tlazy dog." | wordWrap 19 }}

Output:

    The quick brown
    fox jumps over the
    lazy dog.

Template:

    {{ "Γαζέες καὶ μυρτιὲς δὲν θὰ βρῶ πιὰ στὸ χρυσαφὶ ξέφωτο" | wordWrap 19 }}

Output:

    Γαζέες καὶ μυρτιὲς
    δὲν θὰ βρῶ πιὰ στὸ
    χρυσαφὶ ξέφωτο

# Contributing

1. Fork the repository on GitHub
2. Create a named feature branch (i.e. `add-new-feature`)
3. Write your change
4. Submit a Pull Request

# Authors

- William Bailey - [@cowboysfromhell](https://twitter.com/cowboysfromhell) - ([mail@williambailey.org.uk](mailto:mail@williambailey.org.uk))

# License

Licensed under a [MIT license](LICENSE.txt).
