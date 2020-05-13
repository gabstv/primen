#!/bin/sh

pushd $(git rev-parse --show-toplevel)/tests/tau-aseprite/build

go run ../../../cmd/sprite-utils/tau-aseprite/main.go \
 atlas build --template background-importer.json --strict --imageout 'output/bg#.png' --verbose

go run ../../../cmd/sprite-utils/tau-aseprite/main.go \
 atlas build --template player-importer.json --strict --imageout 'output/player#.png' --verbose

popd