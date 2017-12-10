#!/usr/bin/env bash

GOOS=$1
GOARCH=$2

docker run --rm -v "$PWD":/go/src/github.com/vcraescu/holy-internet \
    -w /go/src/github.com/vcraescu/holy-internet \
    -e GOOS=$GOOS \
    -e GOARCH=$GOARCH \
    holyinternet:latest go build -v -o holyinternet-$GOOS-$GOARCH cmd/holyinternet/main.go
