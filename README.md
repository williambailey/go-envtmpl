# envtmpl

Executes text templates using data from environment variables and sends the result to STDOUT.

[![Build Status](https://travis-ci.org/williambailey/go-envtmpl.svg)](https://travis-ci.org/williambailey/go-envtmpl)

## Usage:

#### envtmpl tmplDir tmplName.tmpl

Parse **tmplDir/*.tmpl** and renders **tmplName.tmpl** to STDOUT using
environment variables.

#### envtmpl tmplDir/tmplName.tmpl

Parse **tmplDir/*.tmpl** and renders **tmplName.tmpl** to STDOUT using
environment variables.

#### envtmpl -

Read template from STDIN and render to STDOUT using environment variables.

### Exit codes

* 0 - OK.
* 1 - Usage.
* 2 - Template parse error.
* 3 - Template execution error.

### Template Syntax.

See http://golang.org/pkg/text/template/ for template syntax.

See http://code.google.com/p/re2/wiki/Syntax for regular expression syntax.

## Template Functions

In addition to the [actions](http://golang.org/pkg/text/template/#hdr-Actions)
and [functions](http://golang.org/pkg/text/template/#hdr-Functions) provided by
the core [template engine](http://golang.org/pkg/text/template/#pkg-overview),
envtmpl provides the following functions for use in your templates:

* [base32Decode](#base32decode) - Decodes a base32 string.
* [base32Encode](#base32encode) - Encodes a value to base32.
* [base64Decode](#base64decode) - Decodes a base64 string.
* [base64Encode](#base64encode) - Encodes a value to base64.
* [file](#file) - Read the contents of a file.
* [hash](#hash) - Calculate the hex encoded hash of a string.
* [hexDecode](#hexdecode) - Decodes a hex string.
* [hexEncode](#hexencode) - Encodes a value to hex.
* [include](#include) - Include a template.
* [jsonDecode](#jsondecode) - Decodes a JSON string.
* [jsonEncode](#jsonencode) - Encodes a value to JSON.
* [linePrefix](#lineprefix) - Prefix each line.
* [lower](#lower) - Convert to lower case.
* [regexReplace](#regexreplace) - Replace values using a regular expression.
* [slice](#slice) - Construct a substring from a string.
* [slugify](#slugify) - Transform text to a slugified version.
* [split](#split) - Split in a string substrings using another string.
* [title](#title) - Convert to title case.
* [trimPrefix](#trimprefix) - Remove leading prefix.
* [trimSpace](#trimspace) - Remove all leading and trailing white space.
* [trimSuffix](#trimsuffix) - Remove trailing suffix.
* [upper](#upper) - Convert to upper case.
* [url](#url) - Parse a URL from a string.
* [urlEscape](#urlescape) - Escapes the string so it can be safely placed inside a URL query.
* [urlUnescape](#urlunescape) - Unescapes a URL query string value.
* [uuid](#uuid) - Create a random (v4) UUID.
* [wordWrap](#wordwrap) - Wraps text to a given number of runes.

### base32Decode

Decodes a base32 string.

Template:

    {{ "JBSWY3DPEBLU6USMIQQQ====" | base32Decode }}

Output:

    Hello WORLD!

### base32Encode

Encodes a value to base32.

Template:

    {{ "Hello WORLD!" | base32Encode }}

Output:

    JBSWY3DPEBLU6USMIQQQ====

### base64Decode

Decodes a base64 string.

Template:

    {{ "SGVsbG8gV09STEQh" | base64Decode }}

Output:

    Hello WORLD!

### base64Encode

Encodes a value to base64.

Template:

    {{ "Hello WORLD!" | base64Encode }}

Output:

    SGVsbG8gV09STEQh

### file

Read the contents of a file.

Template:

    {{ file "../example/hello.txt" }}

Output:

    Hello, 世界

### hash

Calculate the hex encoded hash of a string. You can optionally specify a key to
produce a HMAC string. The following hash algorithms are supported: adler32,
crc32, crc64ecma, crc64iso, fnv1-32, fnv1-64, fnv1a-32, fnv1a-64, md5, sha1,
sha224, sha256, sha384, sha512.

Template:

    {{ "Hello World!" | hash "adler32" }}
    {{ "Hello World!" | hash "adler32" "a key" }}

Output:

    1c49043e
    052901bf

Template:

    {{ "Hello World!" | hash "crc32" }}
    {{ "Hello World!" | hash "crc32" "a key" }}

Output:

    1c291ca3
    478193aa

Template:

    {{ "Hello World!" | hash "crc64ecma" }}
    {{ "Hello World!" | hash "crc64ecma" "a key" }}

Output:

    75045245c9ea6fe2
    387589d9abc2595d

Template:

    {{ "Hello World!" | hash "crc64iso" }}
    {{ "Hello World!" | hash "crc64iso" "a key" }}

Output:

    7db9cf17f71cd9ac
    e349c05d90529684

Template:

    {{ "Hello World!" | hash "fnv1-32" }}
    {{ "Hello World!" | hash "fnv1-32" "a key" }}

Output:

    12a9a41c
    44d4e845

Template:

    {{ "Hello World!" | hash "fnv1-64" }}
    {{ "Hello World!" | hash "fnv1-64" "a key" }}

Output:

    8e59dd02f68c387c
    bbf3e38bba595be5

Template:

    {{ "Hello World!" | hash "fnv1a-32" }}
    {{ "Hello World!" | hash "fnv1a-32" "a key" }}

Output:

    b1ea4872
    17ef57fd

Template:

    {{ "Hello World!" | hash "fnv1a-64" }}
    {{ "Hello World!" | hash "fnv1a-64" "a key" }}

Output:

    8c0ec8d1fb9e6e32
    ee9b6cfc8ca44005

Template:

    {{ "Hello World!" | hash "md5" }}
    {{ "Hello World!" | hash "md5" "a key" }}

Output:

    ed076287532e86365e841e92bfc50d8c
    3997a224c5ed2b57cf179a38ec61f455

Template:

    {{ "Hello World!" | hash "sha1" }}
    {{ "Hello World!" | hash "sha1" "a key" }}

Output:

    2ef7bde608ce5404e97d5f042f95f89f1c232871
    edff5c450e50bbc39c684f96f9647a7b2a412c42

Template:

    {{ "Hello World!" | hash "sha224" }}
    {{ "Hello World!" | hash "sha224" "a key" }}

Output:

    4575bb4ec129df6380cedde6d71217fe0536f8ffc4e18bca530a7a1b
    e2ea9d5f84029e65ad0185986f817fa1677e6c9868bbb80aa69b3052

Template:

    {{ "Hello World!" | hash "sha256" }}
    {{ "Hello World!" | hash "sha256" "a key" }}

Output:

    7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069
    31f704a351e00973d2320c7b55eeeab0cd4e53a6b5d281b172596eaf69667cfa

Template:

    {{ "Hello World!" | hash "sha384" }}
    {{ "Hello World!" | hash "sha384" "a key" }}

Output:

    bfd76c0ebbd006fee583410547c1887b0292be76d582d96c242d2a792723e3fd6fd061f9d5cfd13b8f961358e6adba4a
    6a51e0a4ec0cdc6698bccc183dfcd877376c3860871781d705e62286be9fe4c7df14f6c8a6100919e0b84032dd672f5e

Template:

    {{ "Hello World!" | hash "sha512" }}
    {{ "Hello World!" | hash "sha512" "a key" }}

Output:

    861844d6704e8573fec34d967e20bcfef3d424cf48be04e6dc08f2bd58c729743371015ead891cc3cf1c9d34b49264b510751b1ff9e537937bc46b5d6ff4ecc8
    a22b0da52cc34b891eb7575afc4cc73189a0c770a039d6d4c28f3f5a45ea62fce3b73094b63338a32ef5d3c7a20f28331d5ea899c58b8cebdf21aa3469cf3731

### hexDecode

Decodes a hex string.

Template:

    {{ "48656c6c6f20574f524c4421" | hexDecode }}

Output:

    Hello WORLD!

### hexEncode

Encodes a value to hex.

Template:

    {{ "Hello WORLD!" | hexEncode }}

Output:

    48656c6c6f20574f524c4421

### include

Include a template. Differs from the template keyword in that it can accept the
template name as part of a pipeline.

Template:

    {{define "ex"}}FOO is {{ .FOO }}{{end}}{{ $t := "ex" }}>>{{ include $t . }}<<

Output:

    >>FOO is foo<<

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

### slice

Construct a substring from a string.

Template:

    {{ "ᛁᚳ᛫ᛗᚨᚷ᛫ᚷᛚᚨᛋ᛫" | slice 3 7 }}

Output:

    ᛗᚨᚷ᛫

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
    

### urlEscape

Escapes the string so it can be safely placed inside a URL query.

Template:

    {{ "Hello World!" | urlEscape }}

Output:

    Hello+World%21

### urlUnescape

Unescapes a URL query string value.

Template:

    {{ "Hello+World%21" | urlUnescape }}

Output:

    Hello World!

### uuid

Create a random (v4) UUID.

Template:

    {{ uuid }}

Output:

    313d5da8-9b30-426f-b1b7-1c1d914e0da8

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
