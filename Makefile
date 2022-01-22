build:			## Generate application binaries
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o build/folder_manager cmd/main.go || exit $?

test:
	go test ./...
