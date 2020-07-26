#!/bin/bash

go get aletheia.icu/broccoli@v1.0.3
go generate ./...
go run main.go