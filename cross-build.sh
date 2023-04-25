# 64-bit
GOOS=linux GOARCH=amd64 go build -o ./bin/goingo-amd64-linux .
GOOS=windows GOARCH=amd64 go build -o ./bin/goingo-amd64.exe .
