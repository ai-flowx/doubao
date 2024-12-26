#!/bin/bash

ldflags="-s -w"
target="tokenization"

go env -w GOPROXY=https://goproxy.cn,direct

CGO_ENABLED=1 GOARCH=$(go env GOARCH) GOOS=$(go env GOOS) go build -ldflags "$ldflags" -o bin/$target main.go
