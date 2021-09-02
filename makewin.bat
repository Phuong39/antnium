@echo off

IF "%1"=="runserver" (
    go run cmd\server\server.go --listenaddr 127.0.0.1:8080
) ELSE IF "%1"=="runclient" (
    go run cmd\client\client.go
) ELSE IF "%1"=="rundownstreamclient" (
    go run cmd\downstreamclient\downstreamclient.go 
) ELSE IF "%1"=="server" (
    go build cmd\server\server.go 
    SET GOOS=linux
    go build -o server.elf cmd\server\server.go 
) ELSE IF "%1"=="client" (
    go build cmd\client\client.go
) ELSE IF "%1"=="downstreamclient" (
    go build cmd\downstreamclient\downstreamclient.go 
) ELSE IF "%1"=="deploy" (
    .\makewin.bat client
    .\makewin.bat server
    mkdir build\upload
    mkdir build\static
    copy server.elf build\
    copy server.exe build\
    copy client.exe build\static\
) ELSE IF "%1"=="coverage" (
    go test -coverprofile="coverage.out"
    go tool cover -html="coverage.out"
) ELSE (
    echo "Unknown: %1"
)