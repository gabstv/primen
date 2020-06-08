#!/bin/sh

ASEPRIMEN=$(git rev-parse --show-toplevel)/cmd/aseprimen/main.go

echo $ASEPRIMEN

go run $ASEPRIMEN build people.tmpl.json