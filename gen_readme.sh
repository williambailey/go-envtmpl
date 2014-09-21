#!/bin/sh
go build -o ./envtmpl/envtmpl ./envtmpl/*.go
u=`./envtmpl/envtmpl 2>&1`
USAGE=$u ./envtmpl/envtmpl ./example readme.tmpl > README.md
