#!/bin/bash

find . -name \*.go -not -path "./vendor/*" -not -path "./build/*" -exec goimports -w {} \;
