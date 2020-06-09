#!/bin/sh

go build -tags example .

# go tool pprof --pdf ./layers prof-file > file.pdf