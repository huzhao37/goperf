@echo off
@color 06

SET CGO_ENABLED=0

SET GOOS=windows

SET GOARCH=amd64

go build -o ../build/goperf.exe ../cmd/main.go


pause