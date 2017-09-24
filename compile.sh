#!/bin/bash

V=0.4.6

# Windows
GOOS=windows GOARCH=386 go build -o BXGo-$V-Windows-i386.exe ./app
GOOS=windows GOARCH=amd64 go build -o BXGo-$V-Windows-amd64.exe ./app

# macOS
GOOS=darwin GOARCH=amd64 go build -o BXGo-$V-macOS ./app

# Linux
GOOS=linux GOARCH=386 go build -o BXGo-$V-Linux-i386 ./app
GOOS=linux GOARCH=amd64 go build -o BXGo-$V-Linux-amd64 ./app
