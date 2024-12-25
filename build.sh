#!/bin/bash
# By José Puga. 2024. GPL3 License
# Compiles project to Linux & Win, 64bits.

VERSION=$(git describe --tags)
FLAGS="-w -s -X main.version=$VERSION"

for so in linux windows; do
    GOOS=$so go build -o bin/ -ldflags="$FLAGS" .    
done