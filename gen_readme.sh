#!/bin/bash
GLOBIGNORE="*_test.go"
cmd="go run `pwd`/envtmpl/*.go --"
USAGE=`$cmd 2>&1` $cmd `pwd`/example readme.tmpl > README.md
cat README.md
exit $?
