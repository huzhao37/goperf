@echo off
@color 06


SET GOOS=linux

SET GOARCH=amd64

go build -o ../build/goperf ../cmd/main.go


pausego