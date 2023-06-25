#!/bin/sh

GOARCH=amd64 GOOS=linux go build -o build/functions/recently-played cmd/functions/recently-played.go
zip -jrm build/recently-played.zip build/functions/recently-played
