#!/bin/sh
set -e
go build -o /tmp/out "$(dirname "$0")/app/main.go" 
exec /tmp/out "$@"
