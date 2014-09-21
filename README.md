# envtmpl

Quick and dirty utility that executes a text templates using data from environment variables.

## Install

```bash
go get -u github.com/williambailey/go-envtmpl/envtmpl
```

## Usage

```bash
envtmpl example example.tmpl
```

Full usage is available by invoking `envtmpl` without any arguments.

```bash
envtmpl
```
```
Usage: envtmpl tmplDir tmplName

Parse tmplDir/*.tmpl and renders tmplName to stdout using
the data from environmental variables.

See http://golang.org/pkg/text/template/ for template syntax.

Exit codes:

0 OK.
1 Usage.
2 Template parse error.
3 Template execution error.

Template Functions:

lower :: Convert to lower case.
 * {{ "foo BAR bAz" | lower }}
foo bar baz

title :: Convert to title case.
 * {{ "foo BAR bAz" | title }}
Foo BAR BAz

upper :: Convert to upper case.
 * {{ "foo BAR bAz" | upper }}
FOO BAR BAZ

uuid :: Create a random (v4) UUID.
 * {{ uuid }}
fc1f9511-2fa5-457a-9cc0-b448a0bc4c6a

exit status 1
```

## Contributing

1. Fork the repository on GitHub
2. Create a named feature branch (i.e. `add-new-feature`)
3. Write your change
4. Submit a Pull Request

## Authors

- William Bailey - [@cowboysfromhell](https://twitter.com/cowboysfromhell) - ([mail@williambailey.org.uk](mailto:mail@williambailey.org.uk))

## License

Licensed under a [MIT license](LICENSE.txt).
