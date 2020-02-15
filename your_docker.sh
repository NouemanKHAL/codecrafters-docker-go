#!/bin/sh
go build -o /tmp/out "$(dirname "${BASH_SOURCE[0]}")/app/main.go" 
exec /tmp/out "$@"
