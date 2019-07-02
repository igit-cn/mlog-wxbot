#!/usr/bin/env bash
echo 'building...'
GOOS=linux GOARCH=386 go build
