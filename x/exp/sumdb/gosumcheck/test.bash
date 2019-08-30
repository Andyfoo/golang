#!/bin/bash

set -e
go build -o gosumcheck.exe
export GONOSUMDB=*/text # rsc.io/text but not github.com/Andyfoo/golang/x/text
./gosumcheck.exe "$@" -v test.sum
rm -f ./gosumcheck.exe
