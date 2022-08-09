set goos=linux
set goarch=amd64

go build -o psql ..\postgres\main.go
