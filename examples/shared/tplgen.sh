#!/bin/sh

ASEPRIMEN=$(git rev-parse --show-toplevel)/cmd/aseprimen/main.go

echo $ASEPRIMEN

go run $ASEPRIMEN tplgen --sheet people_slices.json \
  --overwrite -f nearest people.tmpl.json