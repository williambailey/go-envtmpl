package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func run(t *testing.T, env []string, args []string, stdin *[]byte) (int, *bytes.Buffer, *bytes.Buffer) {
	inR, inW, err := os.Pipe()
	if err != nil {
		t.Error(err)
		return -1, nil, nil
	}
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Error(err)
		return -1, nil, nil
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Error(err)
		return -1, nil, nil
	}
	if stdin != nil {
		inW.Write(*stdin)
		inW.Close()
	}
	app := new(env, args, inR, outW, errW)
	r := app.main()
	outW.Close()
	errW.Close()
	o := &bytes.Buffer{}
	io.Copy(o, outR)
	outR.Close()
	e := &bytes.Buffer{}
	io.Copy(e, errR)
	errR.Close()
	return r, o, e
}

func TestInvokeWithoutArgsExitsWithUsage(t *testing.T) {
	r, o, e := run(t, []string{}, []string{"me"}, nil)
	if r != exitUsage {
		t.Errorf(
			"Expecting application to terminate with ExitUsage, %d, got %d.",
			exitUsage,
			r,
		)
	}
	if o.Len() != 0 {
		t.Errorf("Expecting stdout len to be 0, got %d", o.Len())
	}
	u := []byte("Usage:")
	if !bytes.HasPrefix(e.Bytes(), u) {
		t.Errorf("Expecting stderr to start `%s` got `%s`", u, e.Bytes()[:len(u)])
	}
}

func TestInvokeWithHelpFlagDisplaysHelpAndExitsWithUsage(t *testing.T) {
	r, o, e := run(t, []string{}, []string{"me", "-h"}, nil)
	if r != exitUsage {
		t.Errorf(
			"Expecting application to terminate with ExitUsage, %d, got %d.",
			exitUsage,
			r,
		)
	}
	if o.Len() != 0 {
		t.Errorf("Expecting stdout len to be 0, got %d", o.Len())
	}
	u := []byte("# Usage:")
	if !bytes.HasPrefix(e.Bytes(), u) {
		t.Errorf("Expecting stderr to start `%s` got `%s`", u, e.Bytes()[:len(u)])
	}
	// Check that we have all available funcs listed.
	for n, fn := range funcMap {
		b := []byte("\n### " + n + "\n")
		if !bytes.Contains(e.Bytes(), b) {
			t.Errorf("Expecting stderr to contain `%s`", b)
		}
		for _, fex := range fn.example(n) {
			var b bytes.Buffer
			for _, v := range strings.SplitAfter(fex, "\n") {
				b.WriteString("    ")
				b.WriteString(v)
			}
			if !bytes.Contains(e.Bytes(), b.Bytes()) {
				t.Errorf("Expecting stderr to contain `%s` example `%s`.", n, b)
			}
		}
	}
}

func TestInvokeWithEmptyTmplDirExitsWithTemplateParseError(t *testing.T) {
	os.Create("foo")
	defer os.Remove("foo")
	r, o, e := run(t, []string{}, []string{"me", "foo", "bar.tmpl"}, nil)
	if r != exitTemplateParseError {
		t.Errorf(
			"Expecting application to terminate with ExitTemplateParseError, %d, got %d.",
			exitTemplateParseError,
			r,
		)
	}
	if o.Len() != 0 {
		t.Errorf("Expecting stdout len to be 0, got %d", o.Len())
	}
	ex := []byte("Template parse error: template: pattern matches no files: `foo/*.tmpl`")
	if !bytes.Equal(bytes.TrimSpace(e.Bytes()), ex) {
		t.Errorf("Expecting stderr to equal `%s` got `%s`", ex, e.Bytes())
	}
}

func TestInvokeWithUnknownTmplDirExitsWithTemplateParseError(t *testing.T) {
	r, o, e := run(t, []string{}, []string{"me", "foo", "bar.tmpl"}, nil)
	if r != exitTemplateParseError {
		t.Errorf(
			"Expecting application to terminate with ExitTemplateParseError, %d, got %d.",
			exitTemplateParseError,
			r,
		)
	}
	if o.Len() != 0 {
		t.Errorf("Expecting stdout len to be 0, got %d", o.Len())
	}
	ex := []byte("Template parse error: template: pattern matches no files: `foo/*.tmpl`")
	if !bytes.Equal(bytes.TrimSpace(e.Bytes()), ex) {
		t.Errorf("Expecting stderr to equal `%s` got `%s`", ex, e.Bytes())
	}
}

func TestInvokeWithUnknownTmplExitsWithTemplateExecutionError(t *testing.T) {
	os.Create("foo.tmpl")
	defer os.Remove("foo.tmpl")
	r, o, e := run(t, []string{}, []string{"me", ".", "bar.tmpl"}, nil)
	if r != exitTemplateExecutionError {
		t.Errorf(
			"Expecting application to terminate with ExitTemplateExecutionError, %d, got %d.",
			exitTemplateExecutionError,
			r,
		)
	}
	if o.Len() != 0 {
		t.Errorf("Expecting stdout len to be 0, got %d", o.Len())
	}
	ex := []byte("Template execution: template: no template \"bar.tmpl\" associated with template \"me [.]\"")
	if !bytes.Equal(bytes.TrimSpace(e.Bytes()), ex) {
		t.Errorf("Expecting stderr to equal `%s` got `%s`", ex, e.Bytes())
	}
}

func TestInvokeWithEmptyTemplate(t *testing.T) {
	os.Create("foo.tmpl")
	defer os.Remove("foo.tmpl")
	r, o, e := run(t, []string{}, []string{"me", ".", "foo.tmpl"}, nil)
	if r != exitOk {
		t.Errorf(
			"Expecting application to terminate with ExitOk, %d, got %d.",
			exitOk,
			r,
		)
	}
	if e.Len() != 0 {
		t.Errorf("Expecting stderr len to be 0, got %d", e.Len())
	}
	if o.Len() != 0 {
		t.Errorf("Expecting stdout len to be 0, got %d", o.Len())
	}
}

func TestInvokeWithSimpleTemplateAndNoEnv(t *testing.T) {
	fo, _ := os.Create("foo.tmpl")
	defer os.Remove("foo.tmpl")
	fo.Write([]byte(`Hello {{.WHAT}}!`))
	fo.Close()
	r, o, e := run(t, []string{}, []string{"me", ".", "foo.tmpl"}, nil)
	if r != exitOk {
		t.Errorf(
			"Expecting application to terminate with ExitOk, %d, got %d.",
			exitOk,
			r,
		)
	}
	if e.Len() != 0 {
		t.Errorf("Expecting stderr len to be 0, got %d", e.Len())
	}
	ex := []byte(`Hello <no value>!`)
	if !bytes.Equal(o.Bytes(), ex) {
		t.Errorf("Expecting stdout to equal `%s` got `%s`", ex, o.Bytes())
	}
}

func TestInvokeWithSimpleTemplateSingleCommandArg(t *testing.T) {
	fo, _ := os.Create("foo.tmpl")
	defer os.Remove("foo.tmpl")
	fo.Write([]byte(`Hello {{.WHAT}}!`))
	fo.Close()
	r, o, e := run(t, []string{"WHAT=World"}, []string{"me", "./foo.tmpl"}, nil)
	if r != exitOk {
		t.Errorf(
			"Expecting application to terminate with ExitOk, %d, got %d.",
			exitOk,
			r,
		)
	}
	if e.Len() != 0 {
		t.Errorf("Expecting stderr len to be 0, got %d", e.Len())
	}
	ex := []byte(`Hello World!`)
	if !bytes.Equal(o.Bytes(), ex) {
		t.Errorf("Expecting stdout to equal `%s` got `%s`", ex, o.Bytes())
	}
}

func TestInvokeWithSimpleTemplateTwoCommandArg(t *testing.T) {
	fo, _ := os.Create("foo.tmpl")
	defer os.Remove("foo.tmpl")
	fo.Write([]byte(`Hello {{.WHAT}}!`))
	fo.Close()
	r, o, e := run(t, []string{"WHAT=World"}, []string{"me", ".", "foo.tmpl"}, nil)
	if r != exitOk {
		t.Errorf(
			"Expecting application to terminate with ExitOk, %d, got %d.",
			exitOk,
			r,
		)
	}
	if e.Len() != 0 {
		t.Errorf("Expecting stderr len to be 0, got %d", e.Len())
	}
	ex := []byte(`Hello World!`)
	if !bytes.Equal(o.Bytes(), ex) {
		t.Errorf("Expecting stdout to equal `%s` got `%s`", ex, o.Bytes())
	}
}

func TestInvokeWithSimpleTemplateFromStdIn(t *testing.T) {
	in := []byte(`Hello {{.WHAT}}!`)
	r, o, e := run(t, []string{"WHAT=World"}, []string{"me", "-"}, &in)
	if r != exitOk {
		t.Errorf(
			"Expecting application to terminate with ExitOk, %d, got %d.",
			exitOk,
			r,
		)
	}
	if e.Len() != 0 {
		t.Errorf("Expecting stderr len to be 0, got %d", e.Len())
		t.Error(e)
	}
	ex := []byte(`Hello World!`)
	if !bytes.Equal(o.Bytes(), ex) {
		t.Errorf("Expecting stdout to equal `%s` got `%s`", ex, o.Bytes())
	}
}

func TestInvokeWithCustomDelimitersAndSimpleTemplateSingleCommandArg(t *testing.T) {
	fo, _ := os.Create("foo.tmpl")
	defer os.Remove("foo.tmpl")
	fo.Write([]byte(`Hello !!.WHAT]]!`))
	fo.Close()
	r, o, e := run(t, []string{"WHAT=World"}, []string{"me", "-dl", "!!", "-dr", "]]", "./foo.tmpl"}, nil)
	if r != exitOk {
		t.Errorf(
			"Expecting application to terminate with ExitOk, %d, got %d.",
			exitOk,
			r,
		)
	}
	if e.Len() != 0 {
		t.Errorf("Expecting stderr len to be 0, got %d", e.Len())
	}
	ex := []byte(`Hello World!`)
	if !bytes.Equal(o.Bytes(), ex) {
		t.Errorf("Expecting stdout to equal `%s` got `%s`", ex, o.Bytes())
	}
}

func TestInvokeWithCustomLeftHandDelimiterAndSimpleTemplateSingleCommandArg(t *testing.T) {
	fo, _ := os.Create("foo.tmpl")
	defer os.Remove("foo.tmpl")
	fo.Write([]byte(`Hello !!.WHAT}}!`))
	fo.Close()
	r, o, e := run(t, []string{"WHAT=World"}, []string{"me", "-dl", "!!", "./foo.tmpl"}, nil)
	if r != exitOk {
		t.Errorf(
			"Expecting application to terminate with ExitOk, %d, got %d.",
			exitOk,
			r,
		)
	}
	if e.Len() != 0 {
		t.Errorf("Expecting stderr len to be 0, got %d", e.Len())
	}
	ex := []byte(`Hello World!`)
	if !bytes.Equal(o.Bytes(), ex) {
		t.Errorf("Expecting stdout to equal `%s` got `%s`", ex, o.Bytes())
	}
}

func TestInvokeWithCustomRightHandDelimiterAndSimpleTemplateSingleCommandArg(t *testing.T) {
	fo, _ := os.Create("foo.tmpl")
	defer os.Remove("foo.tmpl")
	fo.Write([]byte(`Hello {{.WHAT]]!`))
	fo.Close()
	r, o, e := run(t, []string{"WHAT=World"}, []string{"me", "-dr", "]]", "./foo.tmpl"}, nil)
	if r != exitOk {
		t.Errorf(
			"Expecting application to terminate with ExitOk, %d, got %d.",
			exitOk,
			r,
		)
	}
	if e.Len() != 0 {
		t.Errorf("Expecting stderr len to be 0, got %d", e.Len())
	}
	ex := []byte(`Hello World!`)
	if !bytes.Equal(o.Bytes(), ex) {
		t.Errorf("Expecting stdout to equal `%s` got `%s`", ex, o.Bytes())
	}
}
