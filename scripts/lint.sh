#!/bin/sh

if golangci-lint run; then
	exit 0
fi
exit 1
