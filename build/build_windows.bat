@echo off
@color 06

SET CGO_ENABLED=0

SET GOOS=windows

SET GOARCH=amd64

go build -o ../build/receiver.exe ../receiver/main.go

go build -o ../build/send.exe ../send/main.go


pause