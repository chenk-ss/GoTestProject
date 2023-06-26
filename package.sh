#!/bin/bash

set GOARCH=amd64
set GOOS=linux

go build -o ./goTestProject main.go
