#!/usr/bin/env bash

rm -rf releases && mkdir -p releases

# linux
./build.sh linux amd64 && \
    cp holyinternet-linux-amd64 releases/holyinternet && \
    cd releases/ && \
    tar cvfz holyinternet-linux-amd64.tar.gz holyinternet && \
    rm -f holyinternet && \
    cd ..

echo "-------------------------------------------------"
echo "Linux Release: Done"
echo "-------------------------------------------------"

# osx
go build -v -o holyinternet-osx cmd/holyinternet/main.go && \
    cp holyinternet-osx releases/holyinternet && \
    cd releases && \
    tar cvfz holyinternet-osx.tar.gz holyinternet && \
    rm -f holyinternet && \
    cd ..

echo "-------------------------------------------------"
echo "OSX Release: Done"
echo "-------------------------------------------------"
