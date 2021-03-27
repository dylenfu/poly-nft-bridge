#!/bin/bash

#goimports -d $(find . -type f -name '*.go' -not -path "./vendor/*")

find . -name \*.go -not -path "./vendor/*" -not -path "./build/*" -exec goimports -w {} \;
