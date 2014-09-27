#!/bin/bash
GLOBIGNORE="*_test.go"
cd envtmpl
go build .
USAGE=`./envtmpl 2>&1` \
  HELP_USAGE=`./envtmpl -h 2>&1` \
  ./envtmpl ../example readme.tmpl > ../README.md
status=$?
cat ../README.md
rm envtmpl
exit $status
