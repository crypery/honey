set GOROOT=%PROGRAMFILES%\Go
set GOARCH=amd64
set GOOS=windows
set CGO_ENABLED=0
set GO111MODULE=on
go build -ldflags="-s -w" -trimpath -o honey.exe .
