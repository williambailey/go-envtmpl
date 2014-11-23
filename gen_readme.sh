#!/bin/bash
USAGE=`./bin/envtmpl 2>&1` \
  HELP_USAGE=`./bin/envtmpl -h 2>&1` \
    ./bin/envtmpl ./example readme.tmpl > ./README.md
status=$?
cat ./README.md
exit $status
