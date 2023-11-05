#!/bin/sh

GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -o bootstrap -ldflags="-s -w" cmd/functions/recently-played.go
zip -jrm build/recently-played.zip bootstrap
