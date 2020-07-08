@echo off
@color 06


SET GOOS=linux

SET GOARCH=amd64

go build -o ../build/receiver ../receiver/main.go

go build -o ../build/send ../send/main.go


pausego