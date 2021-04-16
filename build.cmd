@echo off
set cmd=go build -i -ldflags "-s -w" -tags "divert_cgo" -v -o bin/PProxy-cli.exe main.go
%cmd%